package context

import (
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/inventory"
)

func (p *ContextPlayer) handleCreativeInventoryAction(packet interface{}) {
	creativeInventoryActionPacket := packet.(*packets.SCreaviteInventoryAction)
	if creativeInventoryActionPacket.ClickedItem.BlockID == -1 {
		p.Inventory.SetSlot(int(creativeInventoryActionPacket.Slot), nil)
	} else if creativeInventoryActionPacket.Slot != -1 {
		p.Inventory.SetSlot(int(creativeInventoryActionPacket.Slot), &inventory.Slot{
			Item:  inventory.Item(creativeInventoryActionPacket.ClickedItem.BlockID),
			State: creativeInventoryActionPacket.ClickedItem.ItemState,
			Count: creativeInventoryActionPacket.ClickedItem.ItemCount,
			NBT:   creativeInventoryActionPacket.ClickedItem.NBT,
		})
	}
}

func (p *ContextPlayer) handleHeldItemChange(packet interface{}) {
	heldItemChangePacket := packet.(*packets.SHeldItemChange)
	p.Mu.Lock()
	p.HeldSlot = heldItemChangePacket.Slot + 36
	p.Mu.Unlock()
}
