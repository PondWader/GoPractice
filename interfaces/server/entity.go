package server_interfaces

import (
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/server/structs"
)

type Entity interface {
	Teleport(*structs.Location)
	Type() string
	SpawnEntityForClient(*protocol.ProtocolClient)
	AddEntityToView(int32, Entity)
	RemoveEntityFromView(int32, bool)
}
