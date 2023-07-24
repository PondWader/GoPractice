package context

import (
	"fmt"

	"github.com/PondWader/GoPractice/protocol"
)

func (p *ContextPlayer) handleBlockPlace(packet interface{}) {
	placePacket := packet.(*protocol.SPlayerBlockPlacement)

	x := placePacket.Location.X
	y := placePacket.Location.Y
	z := placePacket.Location.Z

	if y >= 0 && y <= 255 {
		fmt.Println(int(x&0xf), int(y&0xf), int(z&0xf))
		p.Context.World.GetChunk(x>>4, z>>4).SetBlock(int(x&0xf), int(y&0xf), int(z&0xf), 1)
	} else {
		p.Context.Events.Emit("itemActivated", p)
	}
}
func (p *ContextPlayer) handleDigging(packet interface{}) {

}
