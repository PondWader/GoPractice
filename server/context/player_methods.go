package context

import (
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/server/structs"
	"github.com/PondWader/GoPractice/server/world"
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

func (p *ContextPlayer) streamChunks() {
	p.Mu.Lock()
	centralChunkX := p.Position.GetBlockX() >> 4
	centralChunkZ := p.Position.GetBlockZ() >> 4
	viewDistance := int32(p.Context.config.ViewDistance)

	chunksToBeUnloaded := p.loadedChunks
	p.loadedChunks = make(map[string]*world.ChunkKey)

	for x := (centralChunkX - viewDistance); x <= (centralChunkX + viewDistance); x++ {
		for z := (centralChunkZ - viewDistance); z <= (centralChunkZ + viewDistance); z++ {
			key := world.GetChunkKey(x, z)
			keyStr := key.String()

			if chunksToBeUnloaded[keyStr] != nil {
				delete(chunksToBeUnloaded, keyStr)
				p.loadedChunks[keyStr] = key
				continue
			}

			p.loadedChunks[keyStr] = key
			chunkData := p.Context.World.GetChunkData(x, z)
			p.Client.WritePacket(0x21, protocol.Serialize(chunkData))
		}
	}

	for _, key := range chunksToBeUnloaded {
		chunkData := world.GetEmptyChunk(key.X, key.Z)
		p.Client.WritePacket(0x21, protocol.Serialize(chunkData))
	}

	p.Mu.Unlock()
}
