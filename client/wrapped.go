package client

import (
	"github.com/gocaine/go-dart/hardware"
)

//WrappedClient is structure for manipulate Client
type WrappedClient struct {
	CurrentGameID   int
	LatestGameState GameState
	client          *DartClient
	board           string
}

//NewWrappedClient is WrappedClient constructor
func NewWrappedClient(endpointURL string, board string) *WrappedClient {
	return &WrappedClient{client: NewClient(endpointURL, board), board: board}
}

// Consume is
func (wrapped *WrappedClient) Consume(event hardware.InputEvent) {
	// simply delegate to the client
	wrapped.client.FireDart(event.Sector, event.Multiplier)
}

// Shutdown is
func (wrapped *WrappedClient) Shutdown() {
	wrapped.client.Shutdown()
}
