package context

import (
	"sync"

	"github.com/PondWader/GoPractice/config"
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/server/structs"
	"github.com/PondWader/GoPractice/server/world"
)

// The context package is used to create each "context"
// Each context can essentially be thought of as a running world, for example the lobby is a context and each game running is a context
// The context handles showing players to each other and chunk loading and provides utilities to be built on top of

type Context struct {
	Mu      *sync.RWMutex
	World   *world.World
	players map[int32]*ContextPlayer
	config  *config.ServerConfiguration
}

type ContextPlayer struct {
	EntityId     int32
	Client       *protocol.ProtocolClient
	Position     *structs.Location
	Mu           *sync.Mutex
	loadedChunks map[string]*world.ChunkKey
	Context      *Context

	IsOnGround bool
}

func New(world *world.World, config *config.ServerConfiguration) *Context {
	return &Context{
		Mu:      &sync.RWMutex{},
		players: make(map[int32]*ContextPlayer),
		config:  config,
		World:   world,
	}
}

func (c *Context) AddPlayer(client *protocol.ProtocolClient, entityId int32, mu *sync.Mutex) {
	p := &ContextPlayer{
		EntityId: entityId,
		Client:   client,
		Mu:       mu,
		Context:  c,
	}

	p.Teleport(&structs.Location{
		X: 0,
		Y: 60,
		Z: 0,
	})

	/*mu.Lock()
	chunk := world.NewChunk(0, 0)
	chunk.SetBlock(0, 0, 0, 2)
	chunk.SetBlock(0, 5, 1, 2)
	format := chunk.ToFormat()
	p.Client.WritePacket(0x21, protocol.Serialize(format))
	mu.Unlock()*/
	centralChunk := c.World.GetChunk(0, 0)
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			centralChunk.SetBlock(x, 0, z, 1)
		}
	}
	p.streamChunks()

	c.Mu.Lock()
	c.players[entityId] = p
	c.Mu.Unlock()

	p.loadHandlers()
}

func (p *ContextPlayer) loadHandlers() {
	p.Mu.Lock()
	p.Client.SetPacketHandler(&protocol.SPlayerPositionPacket{}, p.handlePlayerPositionUpdate)
	p.Mu.Unlock()
}

func (c *Context) RemovePlayer(entityId int32) {
	c.Mu.Lock()
	delete(c.players, entityId)
	c.Mu.Lock()
}
