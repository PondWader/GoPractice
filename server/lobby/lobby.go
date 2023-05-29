package lobby

import (
	"github.com/PondWader/GoPractice/config"
	"github.com/PondWader/GoPractice/server/context"
	"github.com/PondWader/GoPractice/server/world"
)

func New(config *config.ServerConfiguration) *context.Context {
	return context.New(world.New("lobby"), config)
}
