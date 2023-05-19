package context

import "github.com/PondWader/GoPractice/protocol"

func (p *ContextPlayer) loadHandlers() {
	p.Mu.Lock()
	p.Mu.Unlock()
}

func (p *ContextPlayer) handlePlayerPositionUpdate(packet interface{}) {
	playerPositionPacket := packet.(*protocol.SPlayerPositionPacket)

	p.Mu.Lock()
	p.Position.SetPos(playerPositionPacket.X, playerPositionPacket.Y, playerPositionPacket.Z)
	p.Mu.Lock()
}
