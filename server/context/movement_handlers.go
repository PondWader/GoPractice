package context

import (
	"github.com/PondWader/GoPractice/protocol"
)

func (p *ContextPlayer) handlePlayerPositionUpdate(packet interface{}) {
	playerPositionPacket := packet.(*protocol.SPlayerPositionPacket)

	p.Mu.Lock()
	oldX := p.Position.X
	oldZ := p.Position.Z

	p.Position.SetPos(playerPositionPacket.X, playerPositionPacket.Y, playerPositionPacket.Z)
	p.Mu.Unlock()

	// Detect if the player has changed chunk if so new chunks need to be loaded
	if uint32(oldX)>>4 != uint32(p.Position.X)>>4 || uint32(oldZ)>>4 != uint32(p.Position.Z)>>4 {
		p.streamChunks()
	}
}
