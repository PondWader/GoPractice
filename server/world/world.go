package world

import (
	"fmt"
	"sync"
	"time"

	"github.com/PondWader/GoPractice/protocol/packets"
)

// The world package handles the loading & saving of world data

type World struct {
	Name   string
	mu     *sync.RWMutex
	chunks map[string]*Chunk
}

func New(name string) *World {
	w := &World{
		Name:   name,
		mu:     &sync.RWMutex{},
		chunks: make(map[string]*Chunk),
	}
	go w.runChunkUnloader()
	return w
}

// Gets a chunk in the world and if it doesn't exist loads a new empty chunk in to memory
func (w *World) GetChunk(x int32, z int32) *Chunk {
	w.mu.RLock()

	key := GetChunkKey(x, z).String()
	if w.chunks[key] == nil {
		w.mu.RUnlock()
		chunk := NewChunk(x, z).SetBlock(0, 0, 0, 0) // We set a block to air so clients accept the chunk as it has data (empty chunks are used to tell the client to unload the chunk)
		w.mu.Lock()
		w.chunks[key] = chunk
		w.mu.Unlock()
		return chunk
	}
	chunk := w.chunks[key]
	w.mu.RUnlock()
	return chunk
}

// Gets a chunk if it exists or if it doesn't returns nil
func (w *World) GetChunkOrNil(x int32, z int32) *Chunk {
	w.mu.RLock()
	key := GetChunkKey(x, z).String()
	chunk := w.chunks[key]
	w.mu.RUnlock()
	return chunk
}

// Gets a chunk in the world in packet format without having to load air chunks
func (w *World) GetChunkData(x int32, z int32) *packets.CChunkData {
	w.mu.RLock()
	key := GetChunkKey(x, z).String()
	if w.chunks[key] == nil {
		chunkData := *AirChunk
		chunkData.ChunkX = x
		chunkData.ChunkZ = z
		w.mu.RUnlock()
		return &chunkData
	}
	chunk := w.chunks[key]
	w.mu.RUnlock()
	return chunk.ToFormat()
}

func GetChunkKey(x int32, z int32) *ChunkKey {
	return &ChunkKey{x, z}
}

// Some default chunks to be used
var EmptyChunk *packets.CChunkData = NewChunk(0, 0).ToFormat()
var AirChunk *packets.CChunkData = NewChunk(0, 0).SetBlock(0, 0, 0, 0).ToFormat()

func GetEmptyChunk(x int32, z int32) *packets.CChunkData {
	chunkData := *EmptyChunk
	chunkData.ChunkX = x
	chunkData.ChunkZ = z
	return &chunkData
}

type ChunkKey struct {
	X int32
	Z int32
}

func (key *ChunkKey) String() string {
	return fmt.Sprint(key.X) + "," + fmt.Sprint(key.Z)
}

// Routinely unloads empty chunks
func (w *World) runChunkUnloader() {
	for {
		time.Sleep(time.Second * 30)

		w.mu.RLock()
		chunks := w.chunks
		w.mu.RUnlock()

		for key, chunk := range chunks {
			chunk.mu.RLock()
			if chunk.IsEmpty && len(chunk.entitiesInChunk) == 0 {
				w.mu.Lock()
				delete(w.chunks, key)
				w.mu.Unlock()
			}
			chunk.mu.RUnlock()
		}
	}
}
