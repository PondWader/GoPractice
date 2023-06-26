package context

import (
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/world"
)

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
			p.Client.WritePacket(packets.CChunkDataId, protocol.Serialize(chunkData))
		}
	}

	for _, key := range chunksToBeUnloaded {
		chunkData := world.GetEmptyChunk(key.X, key.Z)
		p.Client.WritePacket(packets.CChunkDataId, protocol.Serialize(chunkData))
	}

	p.Mu.Unlock()
}
