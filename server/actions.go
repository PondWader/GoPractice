package server

import (
	"crypto/rsa"

	"github.com/PondWader/GoPractice/config"
	"gorm.io/gorm"
)

func (s *Server) GetConfig() config.ServerConfiguration {
	return s.Config
}

func (s *Server) GetKeyPair() *rsa.PrivateKey {
	return s.KeyPair
}

func (s *Server) GetPlayerCount() int {
	return len(s.Players)
}

func (s *Server) GetDatabase() *gorm.DB {
	return s.Db
}
