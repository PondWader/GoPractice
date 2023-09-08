package packets

import "github.com/PondWader/GoPractice/utils/nbt"

type PacketId uint8

type Position struct {
	X int32
	Y int32
	Z int32
}

type Slot struct {
	BlockID   int16                  `type:"Short"`
	ItemCount int8                   `type:"Byte" if:"BlockID" notEquals:"-1"`
	ItemState int16                  `type:"Short" if:"BlockID" notEquals:"-1"`
	NBT       map[string]*nbt.NbtTag `type:"NBT" if:"BlockID" notEquals:"-1"`
}

type ChatComponent struct {
	Text string `json:"text"`
	// TODO: make actually meet the chat spec (https://wiki.vg/Chat)
}

type EntityID struct {
	Id int `type:"VarInt"`
}

// Keep alive packet is same both ways
type KeepAlivePacket struct {
	KeepAliveID int `type:"VarInt"`
}
