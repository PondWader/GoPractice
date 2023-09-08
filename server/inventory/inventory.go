package inventory

import (
	"sync"

	"github.com/PondWader/GoPractice/utils/nbt"
)

type Slot struct {
	Item Item
	NBT  map[string]*nbt.NbtTag
}

type Inventory struct {
	slots []*Slot
	mu    *sync.RWMutex
}

func New(size int) *Inventory {
	inv := &Inventory{
		slots: make([]*Slot, size),
		mu:    &sync.RWMutex{},
	}
	return inv
}

func (inv *Inventory) SetSlot(slot int, data *Slot) {
	inv.mu.Lock()
	inv.slots[slot] = data
	inv.mu.Unlock()
}

func (inv *Inventory) GetSlot(slot int) *Slot {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	return inv.slots[slot]
}
