package protocol

import (
	"github.com/PondWader/GoPractice/utils"
	"github.com/google/uuid"
)

// Keep alive packet is same both ways
type KeepAlivePacket struct {
	KeepAliveID int `type:"VarInt"`
}

/*
	Clientbound packets
*/

type CJoinGamePacket struct {
	EntityID         int32  `type:"Int"`
	GameMode         uint8  `type:"UnsignedByte"`
	Dimension        int8   `type:"Byte"`
	Difficulty       uint8  `type:"UnsignedByte"`
	MaxPlayers       uint8  `type:"UnsignedByte"`
	LevelType        string `type:"String"`
	ReducedDebugInfo bool   `type:"Boolean"`
}

type CChatMessage struct {
	Data     ChatComponent `type:"JSON"`
	Position int8          `type:"Byte"`
}

type CSetPlayerPositionAndLook struct {
	X     float64 `type:"Double"`
	Y     float64 `type:"Double"`
	Z     float64 `type:"Double"`
	Yaw   float32 `type:"Float"`
	Pitch float32 `type:"Float"`
	Flags int8    `type:"Byte"`
}

type CHeldItemChangePacket struct {
	Slot int8 `type:"Byte"`
}

type CPlayerAbilitiesPacket struct {
	Flags        int8    `type:"Byte"`
	FlyingSpeed  float32 `type:"Float"`
	WalkingSpeed float32 `type:"Float"`
}

// Player list item action data

// 0: add player
type PlayerListActionAddPlayer struct {
	UUID               *uuid.UUID                    `type:"UUID"`
	Name               string                        `type:"String"`
	NumberOfProperties int                           `type:"VarInt"`
	Properties         []*PlayerListPlayerProperties `type:"Array"`
	GameMode           int                           `type:"VarInt"`
	Ping               int                           `type:"VarInt"`
	HasDisplayName     bool                          `type:"Boolean"`
	DisplayName        ChatComponent                 `type:"JSON" if:"HasDisplayName"`
}

type PlayerListPlayerProperties struct {
	Name     string `type:"String"`
	Value    string `type:"String"`
	IsSigned bool   `type:"Boolean"`
}

// 1: update gamemode
type PlayerListActionUpdateGamemode struct {
	UUID     *uuid.UUID `type:"UUID"`
	GameMode int        `type:"VarInt"`
}

// 2: update latency
type PlayerListActionUpdateLatency struct {
	UUID *uuid.UUID `type:"UUID"`
	Ping int        `type:"VarInt"`
}

// 3: update display name
type PlayerListActionUpdateDisplayName struct {
	UUID           *uuid.UUID    `type:"UUID"`
	HasDisplayName bool          `type:"Boolean"`
	DisplayName    ChatComponent `type:"JSON" if:"HasDisplayName"`
}

// 4: remove player
type PlayerListActionRemovePlayer struct {
	UUID *uuid.UUID `type:"UUID"`
}

type CPLayerListItemPacket struct {
	Action          int         `type:"VarInt"`
	NumberOfPlayers int         `type:"VarInt"`
	Data            interface{} `type:"Array"`
}

type CPlayerListHeaderAndFooter struct {
	Header ChatComponent `type:"JSON"`
	Footer ChatComponent `type:"JSON"`
}

type CChunkData struct {
	ChunkX             int32  `type:"Int"`
	ChunkZ             int32  `type:"Int"`
	GroundUpContinuous bool   `type:"Boolean"`
	PrimaryBitMask     uint16 `type:"UnsignedShort"`
	Size               int    `type:"VarInt"`
	Data               []byte `type:"ByteArray"`
}

/*
	Serverbound packets
*/

type SPlayerPositionPacket struct {
	X        float64 `type:"Double"`
	Y        float64 `type:"Double"`
	Z        float64 `type:"Double"`
	OnGround bool    `type:"Boolean"`
}

type SChatPacket struct {
	Message string `type:"String"`
}

func (client *ProtocolClient) play() {
	client.state = "play"

	// Client gets returned, main will then create a new server player to wrap around it and then begin the packet listener
}

func (client *ProtocolClient) BeginPacketReader() {
	for {
		packetId, data, err := client.readPacket()
		if err != nil {
			return
		}

		var packetFormat interface{}
		switch packetId {
		case 0x00:
			packetFormat = &KeepAlivePacket{}
		case 0x01:
			packetFormat = &SChatPacket{}
		default:
			utils.Error("Received unrecognized packet of ID", packetId, "from", client.Username)
			continue
		}

		err = client.deserialize(data, packetFormat)
		if err != nil {
			break
		}
		client.HandlePacket(packetFormat)
	}
}
