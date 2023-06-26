package context

import (
	server_interfaces "github.com/PondWader/GoPractice/interfaces/server"
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/structs"
)

func (p *ContextPlayer) Teleport(l *structs.Location) {
	p.Mu.Lock()
	p.Position = l
	p.Client.WritePacket(packets.CSetPlayerPositionAndLookId, protocol.Serialize(&protocol.CSetPlayerPositionAndLook{
		X:     l.X,
		Y:     l.Y,
		Z:     l.Z,
		Yaw:   l.Yaw,
		Pitch: l.Pitch,
		Flags: 0,
	}))
	p.Mu.Unlock()
}

func (p *ContextPlayer) Type() string {
	return "player"
}

func (p *ContextPlayer) SpawnEntityForClient(client *protocol.ProtocolClient) {
	client.WritePacket(packets.CSpawnPlayerId, protocol.Serialize(&protocol.CSpawnPlayerPacket{
		EntityID:    int(p.EntityId),
		UUID:        p.Client.Uuid,
		X:           p.Position.X,
		Y:           p.Position.Y,
		Z:           p.Position.Z,
		Yaw:         p.Position.GetYawAngle(),
		Pitch:       p.Position.GetPitchAngle(),
		CurrentItem: 0,
		Metadata:    []byte{0, 0x08, 0x7F},
	}))
}

func (p *ContextPlayer) AddEntityToView(entityId int32, entity server_interfaces.Entity) {
	p.Mu.Lock()
	p.EntitiesInView[entityId] = entity
	entity.SpawnEntityForClient(p.Client)
	p.Mu.Unlock()
}

func (p *ContextPlayer) RemoveEntityFromView(entityId int32, removalTriggeredBySelf bool) {
	p.Mu.Lock()
	delete(p.EntitiesInView, entityId)

	if removalTriggeredBySelf == false {
		p.Client.WritePacket(packets.CDestroyEntitiesId, protocol.Serialize(&protocol.CDestroyEntitiesPacket{
			Count: 1,
			EntityIDs: []*protocol.EntityID{{
				Id: int(entityId),
			}},
		}))
	}
	p.Mu.Unlock()
}
