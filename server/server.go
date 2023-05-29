package server

/*
	server handles the state of the server, and provides methods for taking actions on all players
*/

import (
	"crypto/rand"
	"crypto/rsa"
	"sync"

	"github.com/PondWader/GoPractice/config"
	"github.com/PondWader/GoPractice/database"
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/server/context"
	"github.com/PondWader/GoPractice/server/lobby"
	"github.com/PondWader/GoPractice/utils"
	"gorm.io/gorm"
)

type Server struct {
	Mu *sync.RWMutex

	KeyPair *rsa.PrivateKey
	Db      *gorm.DB

	Config              config.ServerConfiguration
	Version             string
	Players             map[string]*Player
	entityIdIncrementer int32

	lobby *context.Context
}

func New(cfg config.ServerConfiguration, version string) *Server {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)

	server := &Server{
		Mu:      &sync.RWMutex{},
		KeyPair: key,
		Config:  cfg,
		Db: database.CreateDB(&database.DBConnOptions{
			User:     cfg.DatabaseUser,
			Password: cfg.DatabasePassword,
			Name:     cfg.DatabaseName,
			Host:     cfg.DatabaseHost,
			Port:     cfg.DatabasePort,
		}),
		Version: version,
		lobby:   lobby.New(&cfg),
		Players: map[string]*Player{},
	}

	return server
}

func (s *Server) Broadcast(msg string, position int8) {
	utils.Info(msg)

	s.BroadcastPacket(0x02, protocol.Serialize(&protocol.CChatMessage{
		Data: protocol.ChatComponent{
			Text: msg,
		},
		Position: position,
	}))
}

func (s *Server) BroadcastPacket(packetId int, data []byte) {
	s.Mu.RLock()
	players := s.Players
	s.Mu.RUnlock()

	for _, player := range players {
		// Use goroutine to asynchronously broadcast message so one slow connection doesn't slow everyone receiving it
		go func(p *Player) {
			p.mu.Lock()
			p.client.WritePacket(packetId, data)
			p.mu.Unlock()
		}(player)
	}
}
