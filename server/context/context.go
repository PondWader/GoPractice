package context

import (
	"sync"

	"github.com/PondWader/GoPractice/config"
	server_interfaces "github.com/PondWader/GoPractice/interfaces/server"
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/protocol/enums"
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/structs"
	"github.com/PondWader/GoPractice/server/world"
	"github.com/PondWader/GoPractice/utils"
)

// The context package is used to create each "context"
// Each context can essentially be thought of as a running world, for example the lobby is a context and each game running is a context
// The context handles showing players to each other and chunk loading and provides utilities to be built on top of

type Context struct {
	Mu       *sync.RWMutex
	World    *world.World
	players  map[int32]*ContextPlayer
	config   *config.ServerConfiguration
	building bool
	Events   *utils.EventEmitter
}

type ContextPlayer struct {
	EntityId     int32
	Client       *protocol.ProtocolClient
	Position     *structs.Location
	GameMode     uint8
	Mu           *sync.Mutex
	loadedChunks map[string]*world.ChunkKey
	currentChunk *world.ChunkKey
	Context      *Context

	DisplayedSkinParts uint8

	IsOnGround     bool
	EntitiesInView map[int32]server_interfaces.Entity
}

func New(world *world.World, config *config.ServerConfiguration, building bool) *Context {
	return &Context{
		Mu:       &sync.RWMutex{},
		players:  make(map[int32]*ContextPlayer),
		config:   config,
		World:    world,
		building: building,
		Events:   utils.NewEventEmitter(),
	}
}

func (c *Context) AddPlayer(client *protocol.ProtocolClient, entityId int32, mu *sync.Mutex) {
	p := &ContextPlayer{
		EntityId:       entityId,
		Client:         client,
		Mu:             mu,
		Context:        c,
		Position:       &structs.Location{Y: 60},
		EntitiesInView: make(map[int32]server_interfaces.Entity),
		GameMode:       enums.GamemodeCreative,
	}

	centralChunk := c.World.GetChunk(0, 0)
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			centralChunk.SetBlock(x, 0, z, 1)
		}
	}
	p.streamChunks()
	currentChunk := p.Context.World.GetChunk(p.Position.GetBlockX()>>4, p.Position.GetBlockZ()>>4)
	currentChunk.AddEntity(p.EntityId, p)
	p.currentChunk = currentChunk.GetKey()
	p.updateViewedEntities()

	p.Teleport(&structs.Location{
		X: 0,
		Y: 60,
		Z: 0,
	})

	c.Mu.Lock()
	c.players[entityId] = p
	c.Mu.Unlock()

	p.loadHandlers()
}

func (p *ContextPlayer) loadHandlers() {
	p.Mu.Lock()
	p.Client.SetPacketHandler(&packets.SPlayerPositionPacket{}, p.handlePlayerPositionUpdate)
	p.Client.SetPacketHandler(&packets.SPlayerLookPacket{}, p.handlePlayerLookUpdate)
	p.Client.SetPacketHandler(&packets.SPlayerPositionAndLookPacket{}, p.handlePlayerPositionAndLookUpdate)

	if p.Context.building {
		p.Client.SetPacketHandler(&packets.SPlayerBlockPlacement{}, p.handleBlockPlace)
		p.Client.SetPacketHandler(&packets.SPlayerDigging{}, p.handleDigging)
	}
	p.Mu.Unlock()
}

func (c *Context) RemovePlayer(entityId int32) {
	c.Mu.Lock()
	p := c.players[entityId]
	if p == nil {
		return
	}
	delete(c.players, entityId)
	c.Mu.Unlock()

	p.Mu.Lock()
	if chunk := c.World.GetChunkOrNil(p.currentChunk.X, p.currentChunk.Z); chunk != nil {
		chunk.RemoveEntity(entityId)
	}

	for _, entity := range p.EntitiesInView {
		entity.RemoveEntityFromView(entityId, false)
	}
	p.Mu.Unlock()
}
