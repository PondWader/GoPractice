package world

import (
	"fmt"
	"sync"

	"github.com/PondWader/GoPractice/protocol"
)

type ChunkSection struct {
	blocks [16 * 16 * 16]uint16
}

type Chunk struct {
	mu       *sync.RWMutex
	sections [16]*ChunkSection
	X        int32
	Z        int32
}

func NewChunk(x int32, z int32) *Chunk {
	return &Chunk{
		mu: &sync.RWMutex{},
		X:  x,
		Z:  z,
	}
}

func getBlockIndex(x int, y int, z int) int {
	if x < 0 || z < 0 || x >= 16 || z >= 16 {
		panic("Coords (x=" + fmt.Sprint(x) + ",z=" + fmt.Sprint(z) + ") out of section bounds")
	}
	return ((y & 0xf) << 8) | (z << 4) | x
}

func (c *Chunk) getSection(y int) *ChunkSection {
	c.mu.RLock()
	sectionId := y >> 4
	section := c.sections[sectionId]
	c.mu.RUnlock()
	if section == nil {
		c.mu.Lock()
		section = &ChunkSection{}
		c.sections[sectionId] = section
		c.mu.Unlock()
	}
	return section
}

func (c *Chunk) SetBlock(x int, y int, z int, blockType uint8) *Chunk {
	section := c.getSection(y)

	c.mu.Lock()
	section.blocks[getBlockIndex(x, y, z)] = uint16(blockType << 4)
	c.mu.Unlock()

	return c
}

func (c *Chunk) SetState(x int, y int, z int, state uint8) {
	if state < 0 || state >= 16 {
		panic(fmt.Sprint(state) + " is not a valid state, should be 0-15.")
	}
	section := c.getSection(y)

	c.mu.Lock()
	blockIndex := getBlockIndex(x, y, z)
	currentBlockData := section.blocks[blockIndex]
	section.blocks[blockIndex] = uint16((currentBlockData & 0xfff0) | uint16(state))
	c.mu.Unlock()
}

// Massive thanks to GlowstoneMC for the code this is based off
// https://github.com/GlowstoneMC/Glowstone/blob/d3ed79ea7d284df1d2cd1945bf53d5652962a34f/src/main/java/net/glowstone/GlowChunk.java#L673
func (c *Chunk) ToFormat() *protocol.CChunkData {
	c.mu.RLock()
	chunkData := &protocol.CChunkData{
		GroundUpContinuous: true,
	}

	var bitMask uint16 = 0
	dataSize := 256 // Starting at 256 adds 256 0x00 bytes for the biome
	for i := len(c.sections) - 1; i >= 0; i-- {
		if c.sections[i] == nil {
			bitMask <<= 1
		} else {
			bitMask <<= 1
			bitMask += 1

			BLOCKS_IN_SECTION := 16 * 16 * 16
			dataSize += BLOCKS_IN_SECTION * 5 / 2
			dataSize += BLOCKS_IN_SECTION / 2 // Space for skylight
		}
	}
	chunkData.PrimaryBitMask = bitMask
	chunkData.Size = dataSize

	chunkData.ChunkX = c.X
	chunkData.ChunkZ = c.Z

	data := make([]byte, dataSize)
	pos := 0

	for _, section := range c.sections {
		if section == nil {
			continue
		}
		//Block data
		for _, block := range section.blocks {
			// Write the uint16 value as big endian
			data[pos] = byte(block & 0xff)
			data[pos+1] = byte(block >> 8)
			pos += 2
		}

		// Set light level to 10 for all blocks & sky light for now
		for i := 0; i < 16*16*16; i++ {
			data[pos+i] = 0xaa // 10 in each 4 bit segment
		}
	}

	chunkData.Data = data

	c.mu.RUnlock()
	return chunkData
}

func (c *Chunk) ToSaveFormat(world string) []byte {
	format := c.ToFormat()

	blockDataSize := 0
	c.mu.RLock()
	for _, section := range c.sections {
		if section == nil {
			continue
		}
		blockDataSize += 16 * 16 * 16 * 2
	}
	c.mu.RUnlock()

	data := append([]byte{uint8(format.PrimaryBitMask), uint8(format.PrimaryBitMask >> 8)}, format.Data[:blockDataSize]...)
	return data
}

func ChunkFromSave(x int32, z int32, data []byte) *Chunk {
	chunk := NewChunk(x, z)

	bitMask := uint16(data[0]) + (uint16(data[1]) << 8)
	var currentChunkBitMask uint16 = 1

	offset := 2
	for i := 0; i < 16; i++ {
		if bitMask&currentChunkBitMask == currentChunkBitMask {
			section := &ChunkSection{}
			chunk.sections[i] = section
			for j := 0; j < 16*16*16; j++ {
				section.blocks[j] = uint16(data[offset]) + (uint16(data[offset+1]) << 8)
				offset += 2
			}
		}
		currentChunkBitMask <<= 1
	}

	return chunk
}

func (c *Chunk) getKey() string {
	return fmt.Sprint(c.X) + "," + fmt.Sprint(c.Z)
}
