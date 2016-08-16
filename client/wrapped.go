package client

import (
	"github.com/gocaine/go-dart/hardware"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

//WrappedClient is structure for manipulate Client
type WrappedClient struct {
	CurrentGameID   int
	LatestGameState GameState
	client          *DartClient
}

//NewWrappedClient is WrappedClient constructor
func NewWrappedClient(endpointURL string) *WrappedClient {
	return &WrappedClient{client: NewClient(endpointURL)}
}

// Consume is
func (wrapped *WrappedClient) Consume(event hardware.InputEvent) {
	// simply delegate to the client
	board := viper.GetString("board")
	log.Printf("Board: %s", board)
	wrapped.client.FireDart(board, event.Sector, event.Multiplier)
}

// Shutdown is
func (wrapped *WrappedClient) Shutdown() {
	log.Info("Mocked client is shuting down")
}
