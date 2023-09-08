package context

import (
	"fmt"

	"github.com/PondWader/GoPractice/protocol/packets"
)

func (p *ContextPlayer) handleCreativeInventoryAction(packet interface{}) {
	creativeInventoryActionPacket := packet.(*packets.SCreaviteInventoryAction)
	fmt.Println(creativeInventoryActionPacket, creativeInventoryActionPacket.ClickedItem)
}
