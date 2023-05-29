package world

import (
	"fmt"
	"sync"

	"github.com/PondWader/GoPractice/protocol"
)

// The world package handles the loading & saving of world data

type World struct {
	Name   string
	mu     *sync.RWMutex
	chunks map[string]*Chunk
}

func New(name string) *World {
	return &World{
		Name:   name,
		mu:     &sync.RWMutex{},
		chunks: make(map[string]*Chunk),
	}
}

// Gets a chunk in the world and if it doesn't exist loads a new empty chunk in to memory
func (w *World) GetChunk(x int32, z int32) *Chunk {
	w.mu.RLock()
	key := GetChunkKey(x, z).String()
	if w.chunks[key] == nil {
		w.mu.RUnlock()
		chunk := NewChunk(0, 0).SetBlock(0, 0, 0, 0) // We set a block to air so clients accept the chunk as it has data (empty chunks are used to tell the client to unload the chunk)
		w.mu.Lock()
		w.chunks[key] = chunk
		w.mu.Unlock()
		return chunk
	}
	chunk := w.chunks[key]
	w.mu.RUnlock()
	return chunk
}

// Gets a chunk in the world in packet format without having to load air chunks
func (w *World) GetChunkData(x int32, z int32) *protocol.CChunkData {
	w.mu.RLock()
	key := GetChunkKey(x, z).String()
	if w.chunks[key] == nil {
		chunkData := *AirChunk
		chunkData.ChunkX = x
		chunkData.ChunkZ = z
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
var EmptyChunk *protocol.CChunkData = NewChunk(0, 0).ToFormat()
var AirChunk *protocol.CChunkData = NewChunk(0, 0).SetBlock(0, 0, 0, 0).ToFormat()

func GetEmptyChunk(x int32, z int32) *protocol.CChunkData {
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
