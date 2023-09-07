package context

import (
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/world"
)

func getBlockLocation(pos *packets.Position, face int8) (x int32, y int32, z int32) {
	x, y, z = pos.X, pos.Y, pos.Z
	switch face {
	case 0:
		y -= 1
	case 1:
		y += 1
	case 2:
		z -= 1
	case 3:
		z += 1
	case 4:
		x -= 1
	case 5:
		x += 1
	}
	return
}

func (p *ContextPlayer) handleBlockPlace(packet interface{}) {
	placePacket := packet.(*packets.SPlayerBlockPlacement)

	x, y, z := getBlockLocation(placePacket.Location, placePacket.Face)

	if y >= 0 && y <= 255 {
		xInChunk, yInChunk, zInChunk := world.CoordsInChunk(int(x&0xf), int(y&0xf), int(z&0xf))
		p.Context.World.GetChunk(x>>4, z>>4).SetBlock(xInChunk, yInChunk, zInChunk, 1)
	} else {
		p.Context.Events.Emit("itemActivated", p)
	}
}
func (p *ContextPlayer) handleDigging(packet interface{}) {

}
