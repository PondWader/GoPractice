package server

import (
	"math/rand"
	"sync"
	"time"

	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/protocol/packets"
	"github.com/PondWader/GoPractice/server/context"
)

type Player struct {
	mu     *sync.Mutex
	client *protocol.ProtocolClient
	server *Server

	lastSentKeepAlive     int
	lastReceivedKeepAlive int

	entityId       int32
	currentContext *context.Context
}

func NewPlayer(client *protocol.ProtocolClient, server *Server) *Player {
	server.Mu.RLock()
	players := server.Players
	server.Mu.RUnlock()
	if players[client.Uuid.String()] != nil {
		client.Disconnect("§cUser is already connected from another location.")
		return nil
	}

	server.Mu.Lock()
	p := &Player{
		client:   client,
		server:   server,
		entityId: server.entityIdIncrementer,
		mu:       &sync.Mutex{},
	}
	server.entityIdIncrementer++
	server.Players[client.Uuid.String()] = p
	p.currentContext = server.lobby
	server.Mu.Unlock()

	go p.keepAlive()

	p.loadInPlayer()

	client.SetEndHandler(p.handleEnd)

	client.SetPacketHandler(&protocol.KeepAlivePacket{}, p.handleKeepAlive)
	client.SetPacketHandler(&protocol.SChatPacket{}, p.handleChat)
	client.BeginPacketReader()

	return p
}

func (p *Player) handleEnd() {
	p.removeFromPlayerlist()

	p.server.Mu.Lock()
	delete(p.server.Players, p.client.Uuid.String())
	p.currentContext.RemovePlayer(p.entityId)
	p.server.Mu.Unlock()

	p.server.Broadcast("§e"+p.client.Username+" has left the server.", 1)

}

func (p *Player) loadInPlayer() {
	p.mu.Lock()

	p.client.WritePacket(packets.CJoinGameId, protocol.Serialize(&protocol.CJoinGamePacket{
		EntityID:         p.entityId,
		GameMode:         1, //Temp: creative // Adventure
		Dimension:        0,
		Difficulty:       2, // Normal difficulty
		MaxPlayers:       100,
		LevelType:        "flat",
		ReducedDebugInfo: false,
	}))

	p.client.WritePacket(packets.CPlayAbilitiesId, protocol.Serialize(&protocol.CPlayerAbilitiesPacket{
		Flags:        0x04,
		FlyingSpeed:  0.05,
		WalkingSpeed: 0.1,
	}))

	p.client.WritePacket(packets.CHeldItemChangeId, protocol.Serialize(&protocol.CHeldItemChangePacket{
		Slot: 0,
	}))

	p.mu.Unlock()

	p.loadPlayerList()
	p.addToPlayerlist()

	p.server.Broadcast("§e"+p.client.Username+" has joined the server.", 1)
	p.currentContext.AddPlayer(p.client, p.entityId, p.mu)
}

// Function that runs the goroutine responsible for sending keep alives and checking if the client has timed out
func (p *Player) keepAlive() {
	for {
		time.Sleep(time.Second * 15)

		p.mu.Lock()
		if p.client.Ended {
			break
		}

		if p.lastReceivedKeepAlive != p.lastSentKeepAlive {
			p.client.Disconnect("Timed out.")
			break
		}

		p.lastSentKeepAlive = int(rand.Int31())
		p.client.WritePacket(packets.CKeepAliveId, protocol.Serialize(&protocol.KeepAlivePacket{
			KeepAliveID: p.lastSentKeepAlive,
		}))
		p.mu.Unlock()
	}
}
