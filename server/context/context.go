package context

import (
	"sync"

	"github.com/PondWader/GoPractice/protocol"
)

// The context package is used to create each "context"
// Each context can essentially be thought of as a running world, for example the lobby is a context and each game running is a context
// The context handles showing players to each other and chunk loading and provides utilities to be built on top of

type Context struct {
	players map[int32]*ContextPlayer
}

type ContextPlayer struct {
	EntityId int32
	Client   *protocol.ProtocolClient
	mu       *sync.Mutex
}

func (c *Context) AddPlayer(client *protocol.ProtocolClient, entityId int32, mu *sync.Mutex) {
	p := &ContextPlayer{
		EntityId: entityId,
		Client:   client,
		mu:       mu,
	}

	c.players[entityId] = p
}

func (c *Context) RemovePlayer(entityId int32) {
	delete(c.players, entityId)
}
