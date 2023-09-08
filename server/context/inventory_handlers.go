package context

import (
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/inventory"
)

func (p *ContextPlayer) handleCreativeInventoryAction(packet interface{}) {
	creativeInventoryActionPacket := packet.(*packets.SCreaviteInventoryAction)
	if creativeInventoryActionPacket.ClickedItem.BlockID == -1 {
		p.Inventory.SetSlot(int(creativeInventoryActionPacket.Slot), nil)
	} else {
		p.Inventory.SetSlot(int(creativeInventoryActionPacket.Slot), &inventory.Slot{
			Item: inventory.Item(creativeInventoryActionPacket.ClickedItem.BlockID),
			NBT:  creativeInventoryActionPacket.ClickedItem.NBT,
		})
	}
}
