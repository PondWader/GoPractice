package protocol

import (
	"fmt"
)

type HandshakePacket struct {
	ProtocolVersion int    `type:"VarInt"`
	ServerAddress   string `type:"String"`
	ServerPort      uint16 `type:"UnsignedShort"`
	NextState       int    `type:"VarInt"`
}

func (client *ProtocolClient) handshake() {
	client.state = "handshaking"

	packetId, bytes, err := client.readPacket()
	if err != nil {
		return
	}
	if packetId != 0 {
		client.Disconnect("Received packet ID " + fmt.Sprint(packetId) + " when expecting handshake.")
		return
	}

	handshake := HandshakePacket{}
	err = client.deserialize(bytes, &handshake)
	if err != nil {
		return
	}
	client.protocolVersion = handshake.ProtocolVersion

	if handshake.NextState == 1 {
		client.status()
	} else {
		if handshake.ProtocolVersion != 47 {
			client.Disconnect("This server only supports 1.8.x versions.")
			return
		}
		client.login()
	}
}
