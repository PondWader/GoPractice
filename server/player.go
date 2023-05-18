package server

import (
	"math/rand"
	"sync"
	"time"

	"github.com/PondWader/GoPractice/protocol"
)

type Player struct {
	mu     *sync.Mutex
	client *protocol.ProtocolClient
	server *Server

	lastSentKeepAlive     int
	lastReceivedKeepAlive int

	entityId int32
}

func NewPlayer(client *protocol.ProtocolClient, server *Server) *Player {
	server.Mu.RLock()
	players := server.Players
	server.Mu.RUnlock()
	for _, player := range players {
		if player.client.Username == client.Username {
			client.Disconnect("§cUser is already connected from another location.")
			return nil
		}
	}

	server.Mu.RLock()
	p := &Player{
		client:   client,
		server:   server,
		entityId: server.entityIdIncrementer,
		mu:       &sync.Mutex{},
	}
	server.Mu.RUnlock()
	server.Mu.Lock()
	server.entityIdIncrementer++
	server.Players = append(server.Players, p)
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
}

func (p *Player) loadInPlayer() {
	p.mu.Lock()

	p.client.WritePacket(0x01, protocol.Serialize(&protocol.CJoinGamePacket{
		EntityID:         p.entityId,
		GameMode:         2, // Adventure
		Dimension:        0,
		Difficulty:       2, // Normal difficulty
		MaxPlayers:       100,
		LevelType:        "flat",
		ReducedDebugInfo: false,
	}))

	p.client.WritePacket(0x39, protocol.Serialize(&protocol.CPlayerAbilitiesPacket{
		Flags:        0,
		FlyingSpeed:  0.05,
		WalkingSpeed: 0.1,
	}))

	p.client.WritePacket(0x09, protocol.Serialize(&protocol.CHeldItemChangePacket{
		Slot: 0,
	}))

	// Temp just to test stuff, this should be handled by context
	p.client.WritePacket(0x08, protocol.Serialize(&protocol.CSetPlayerPositionAndLook{
		X:     0,
		Y:     100,
		Z:     0,
		Yaw:   0,
		Pitch: 0,
		Flags: 0,
	}))

	p.mu.Unlock()

	p.loadPlayerList()
	p.addToPlayerlist()
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
		p.client.WritePacket(0x00, protocol.Serialize(&protocol.KeepAlivePacket{
			KeepAliveID: p.lastSentKeepAlive,
		}))
		p.mu.Unlock()
	}
}
