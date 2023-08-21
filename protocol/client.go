package protocol

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"errors"
	"io"
	"net"
	"reflect"
	"time"

	"github.com/PondWader/GoPractice/interfaces"
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/utils"
	"github.com/google/uuid"
)

type ProtocolClient struct {
	conn   net.Conn
	server interfaces.Server

	Ended       bool
	state       string
	compression bool
	encryption  bool

	Username        string
	Skin            string
	Uuid            *uuid.UUID
	protocolVersion int
	decrypter       cipher.Stream
	encrypter       cipher.Stream

	listeners  map[string]func(interface{})
	endHandler func()
}

func NewClient(conn net.Conn, server interfaces.Server) *ProtocolClient {
	client := &ProtocolClient{conn: conn, server: server, listeners: make(map[string]func(interface{}))}
	client.handshake()
	return client
}

type ChatComponent struct {
	Text string `json:"text"`
	// TO DO: make actually meet the chat spec (https://wiki.vg/Chat)
}
type DisconnectPacket struct {
	Reason ChatComponent `type:"JSON"`
}

func (client *ProtocolClient) Disconnect(reason string) {
	if client.Ended {
		return
	}

	client.Ended = true

	if client.state != "play" {
		// Switch state to play so can send disconnect message
		client.WritePacket(0x02, Serialize(&LoginSuccessPacket{
			UUID:     "00000000-0000-0000-0000-000000000000",
			Username: "Player",
		}))
	}
	client.WritePacket(0x40, Serialize(&DisconnectPacket{
		Reason: ChatComponent{
			Text: reason,
		},
	}))

	client.conn.Close()
	if client.Username == "" {
		utils.Info("Connection from", client.conn.RemoteAddr(), "ended. Reason:", reason)
	} else {
		utils.Info(client.Username+"'s", "connection from", client.conn.RemoteAddr(), "ended. Reason:", reason)
	}

	if client.endHandler != nil {
		client.endHandler()
	}
}

func (client *ProtocolClient) WritePacket(packetId packets.PacketId, data []byte) error {
	packetIdVarInt := writeVarInt(int(packetId))
	packetLengthVarInt := writeVarInt(len(packetIdVarInt) + len(data))

	packetData := append(packetIdVarInt, data...)

	if client.compression {
		if len(packetData) > client.server.GetConfig().CompressionThreshold {
			var in bytes.Buffer
			w := zlib.NewWriter(&in)
			w.Write(packetData)
			w.Close()
			data = append(packetLengthVarInt, in.Bytes()...)
			data = append(writeVarInt(len(data)), data...)
		} else {
			// Adding a length of 0 for compression size as the data is not compressed
			data = append(append(writeVarInt(1+len(packetData)), 0), packetData...)
		}
	} else {
		data = append(packetLengthVarInt, packetData...)
	}

	client.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	err := client.writeBytes(data)
	if err != nil {
		client.Disconnect(err.Error())
		return err
	}
	return nil
}

// Reads the packet ID and the decompressed & unencrypted data
// This data can then be passed to the deserialize method with the packet structure to parse it
func (client *ProtocolClient) readPacket() (packets.PacketId, []byte, error) {
	deserializer := client.newDeserializer(client.readBytes)

	dataLength, _, err := deserializer.readVarInt()
	if err != nil {
		client.Disconnect(err.Error())
		return 0, nil, err
	}

	if dataLength <= 0 || dataLength > 10_000 {
		client.Disconnect("Invalid data length")
		return 0, nil, errors.New("Invalid data length")
	}

	if client.compression {
		decompressedLength, bytesRead, err := deserializer.readVarInt()
		if err != nil {
			client.Disconnect(err.Error())
			return 0, nil, err
		}
		dataLength -= bytesRead

		if decompressedLength > 0 {
			compressedData, err := deserializer.readBytes(dataLength)
			if err != nil {
				client.Disconnect(err.Error())
				return 0, nil, err
			}

			var out bytes.Buffer
			r, err := zlib.NewReader(bytes.NewBuffer(compressedData))
			if err != nil {
				client.Disconnect(err.Error())
				return 0, nil, err
			}
			io.Copy(&out, r)
			data := out.Bytes()

			offset := 0
			compressedPacketDeserializer := client.newDeserializer(func(l int) ([]byte, error) {
				if len(data) < offset+l {
					return nil, errors.New("EOF")
				}
				d := data[offset : offset+l]
				offset += l
				return d, nil
			})
			packetId, bytesRead, err := compressedPacketDeserializer.readVarInt()
			if err != nil {
				client.Disconnect(err.Error())
				return 0, nil, err
			}
			decompressedLength -= bytesRead
			if dataLength == 0 {
				return packets.PacketId(packetId), []byte{}, nil
			}

			return packets.PacketId(packetId), data[bytesRead:], nil
		}
	}

	packetId, bytesRead, err := deserializer.readVarInt()
	if err != nil {
		client.Disconnect(err.Error())
		return 0, nil, err
	}
	dataLength -= bytesRead
	if dataLength == 0 {
		return packets.PacketId(packetId), []byte{}, nil
	}

	packetData, err := client.readBytes(dataLength)
	if err != nil {
		client.Disconnect(err.Error())
		return 0, nil, err
	}
	return packets.PacketId(packetId), packetData, nil
}

func (client *ProtocolClient) readBytes(amount int) ([]byte, error) {
	client.conn.SetReadDeadline(time.Now().Add(time.Second * 15))
	bytes := make([]byte, amount)
	_, err := client.conn.Read(bytes)
	if err != nil {
		return []byte{}, err
	}
	if client.encryption == true {
		client.decrypter.XORKeyStream(bytes, bytes)
	}
	return bytes, nil
}

func (client *ProtocolClient) writeBytes(data []byte) error {
	if client.encryption == true {
		client.encrypter.XORKeyStream(data, data)
	}
	_, err := client.conn.Write(data)
	return err
}

func (client *ProtocolClient) SetPacketHandler(packet interface{}, handler func(interface{})) {
	if handler != nil {
		client.listeners[reflect.Indirect(reflect.ValueOf(packet)).Type().Name()] = handler
	} else {
		delete(client.listeners, reflect.TypeOf(packet).Name())
	}
}

func (client *ProtocolClient) HandlePacket(packet interface{}) {
	handler := client.listeners[reflect.Indirect(reflect.ValueOf(packet)).Type().Name()]
	if handler != nil {
		handler(packet)
	}
}

func (client *ProtocolClient) SetEndHandler(handler func()) {
	client.endHandler = handler
}
