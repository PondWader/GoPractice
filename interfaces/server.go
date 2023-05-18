package interfaces

import (
	"crypto/rsa"

	"github.com/PondWader/GoPractice/config"
	"gorm.io/gorm"
)

type Server interface {
	GetConfig() config.ServerConfiguration
	GetPlayerCount() int
	GetKeyPair() *rsa.PrivateKey
	GetDatabase() *gorm.DB
}
