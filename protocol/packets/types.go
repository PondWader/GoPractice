package packets

import "github.com/PondWader/GoPractice/utils/nbt"

type PacketId uint8

type Position struct {
	X int32
	Y int32
	Z int32
}

type Slot struct {
	Present   bool                   `type:"Boolean"`
	ItemID    int                    `type:"VarInt" if:"Present"`
	ItemCount int8                   `type:"Byte" if:"Present"`
	NBT       map[string]*nbt.NbtTag `type:"NBT" if:"Present"`
}

type ChatComponent struct {
	Text string `json:"text"`
	// TO DO: make actually meet the chat spec (https://wiki.vg/Chat)
}

type EntityID struct {
	Id int `type:"VarInt"`
}

// Keep alive packet is same both ways
type KeepAlivePacket struct {
	KeepAliveID int `type:"VarInt"`
}
