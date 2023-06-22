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
		fmt.Println("Block placed at", x, y, z)

		p.Context.World.GetChunk(x>>4, z>>4).SetBlock(int(x), int(y), int(z), 1)
	} else {
		fmt.Println("Item activated", y)
	}

}
