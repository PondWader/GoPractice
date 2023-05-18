package lobby

import (
	"sync"

	"github.com/PondWader/GoPractice/protocol"
)

type Lobby struct {
	players []*protocol.ProtocolClient
}

func (l *Lobby) AddToLobby(client *protocol.ProtocolClient, mu *sync.Mutex) {
	l.players = append(l.players, client)

	newHandler(client, mu)
}
