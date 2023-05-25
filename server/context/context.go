package context

import (
	"sync"

	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/server/structs"
	"github.com/PondWader/GoPractice/server/world"
)

// The context package is used to create each "context"
// Each context can essentially be thought of as a running world, for example the lobby is a context and each game running is a context
// The context handles showing players to each other and chunk loading and provides utilities to be built on top of

type Context struct {
	Mu      *sync.RWMutex
	players map[int32]*ContextPlayer
}

type ContextPlayer struct {
	EntityId int32
	Client   *protocol.ProtocolClient
	Position *structs.Location
	Mu       *sync.Mutex
}

func New() *Context {
	return &Context{
		Mu:      &sync.RWMutex{},
		players: make(map[int32]*ContextPlayer),
	}
}

func (c *Context) AddPlayer(client *protocol.ProtocolClient, entityId int32, mu *sync.Mutex) {
	p := &ContextPlayer{
		EntityId: entityId,
		Client:   client,
		Mu:       mu,
	}

	p.Teleport(&structs.Location{
		X: 0,
		Y: 60,
		Z: 0,
	})

	mu.Lock()
	chunk := world.NewChunk()
	chunk.SetBlock(0, 0, 0, 2)
	chunk.SetBlock(0, 5, 1, 2)
	format := chunk.ToFormat()
	p.Client.WritePacket(0x21, protocol.Serialize(format))

	/*p.Client.WritePacket(0x21, protocol.Serialize(&protocol.CChunkData{
		ChunkX:             0,
		ChunkZ:             0,
		GroundUpContinuous: true,
		PrimaryBitMask:     1,
		Size:               16 * 16 * 16,
		Data:               make([]byte, 16*16*16),
	}))*/
	mu.Unlock()

	c.Mu.Lock()
	c.players[entityId] = p
	c.Mu.Unlock()
}

func (c *Context) RemovePlayer(entityId int32) {
	delete(c.players, entityId)
}
