package lobby

import (
	"sync"

	"github.com/PondWader/GoPractice/protocol"
)

type lobbyHandler struct {
	client *protocol.ProtocolClient
	mu     *sync.Mutex
}

func newHandler(client *protocol.ProtocolClient, mu *sync.Mutex) *lobbyHandler {
	handler := &lobbyHandler{client: client}
	return handler
}
