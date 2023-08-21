package server

import (
	"github.com/PondWader/GoPractice/protocol/packets"
)

func (p *Player) handleKeepAlive(packet interface{}) {
	keepAlivePacket := packet.(*packets.KeepAlivePacket)

	p.mu.Lock()
	p.lastReceivedKeepAlive = keepAlivePacket.KeepAliveID
	p.mu.Unlock()
}

func (p *Player) handleChat(packet interface{}) {
	chatPacket := packet.(*packets.SChatPacket)

	if len(chatPacket.Message) > 100 {
		p.client.Disconnect("Chat message too long (> 100)")
		return
	}

	p.server.Broadcast("§7"+p.client.Username+": §r"+chatPacket.Message, 0)
}
