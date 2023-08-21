package context

import (
	server_interfaces "github.com/PondWader/GoPractice/interfaces/server"
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/protocol/packets"
)

func (p *ContextPlayer) updateViewedEntities() {
	var entityViewDistance int32 = 6
	if int32(p.Context.config.ViewDistance) < entityViewDistance {
		entityViewDistance = int32(p.Context.config.ViewDistance)
	}

	p.Mu.Lock()

	centralChunkX := p.Position.GetBlockX() >> 4
	centralChunkZ := p.Position.GetBlockZ() >> 4

	entitiesToRemove := p.EntitiesInView
	p.EntitiesInView = make(map[int32]server_interfaces.Entity)

	for x := (centralChunkX - entityViewDistance); x <= (centralChunkX + entityViewDistance); x++ {
		for z := (centralChunkZ - entityViewDistance); z <= (centralChunkZ + entityViewDistance); z++ {
			chunk := p.Context.World.GetChunkOrNil(x, z)
			if chunk == nil {
				continue
			}

			for entityId, entity := range chunk.GetEntities() {
				if entityId == p.EntityId {
					continue
				}

				p.EntitiesInView[entityId] = entity

				if entitiesToRemove[entityId] == nil {
					p.Mu.Unlock()
					go entity.AddEntityToView(p.EntityId, p)
					p.AddEntityToView(entityId, entity)
					p.Mu.Lock()
				} else {
					delete(entitiesToRemove, entityId)
				}
			}
		}
	}

	if len(entitiesToRemove) > 0 {
		destroyEntitiesPacket := &packets.CDestroyEntitiesPacket{
			Count:     len(entitiesToRemove),
			EntityIDs: make([]*packets.EntityID, len(entitiesToRemove)),
		}

		i := 0
		for entityId, entity := range entitiesToRemove {
			destroyEntitiesPacket.EntityIDs[i] = &packets.EntityID{Id: int(entityId)}

			p.Mu.Unlock()
			p.RemoveEntityFromView(entityId, true)
			p.Mu.Lock()

			go entity.RemoveEntityFromView(entityId, false)

			i++
		}

		p.Client.WritePacket(packets.CDestroyEntitiesId, protocol.Serialize(destroyEntitiesPacket))
	}

	p.Mu.Unlock()
}
