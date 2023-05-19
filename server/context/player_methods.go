package context

import (
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/server/structs"
)

func (p *ContextPlayer) Teleport(l *structs.Location) {
	p.Mu.Lock()
	p.Position = l
	p.Client.WritePacket(0x08, protocol.Serialize(&protocol.CSetPlayerPositionAndLook{
		X:     l.X,
		Y:     l.Y,
		Z:     l.Z,
		Yaw:   l.Yaw,
		Pitch: l.Pitch,
		Flags: 0,
	}))
	p.Mu.Unlock()
}
