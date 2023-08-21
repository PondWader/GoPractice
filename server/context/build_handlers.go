package context

import (
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/world"
)

func (p *ContextPlayer) handleBlockPlace(packet interface{}) {
	placePacket := packet.(*packets.SPlayerBlockPlacement)

	x := placePacket.Location.X
	y := placePacket.Location.Y
	z := placePacket.Location.Z

	if y >= 0 && y <= 255 {
		xInChunk, yInChunk, zInChunk := world.CoordsInChunk(int(x&0xf), int(y&0xf), int(z&0xf))
		p.Context.World.GetChunk(x>>4, z>>4).SetBlock(xInChunk, yInChunk, zInChunk, 1)
	} else {
		p.Context.Events.Emit("itemActivated", p)
	}
}
func (p *ContextPlayer) handleDigging(packet interface{}) {

}
