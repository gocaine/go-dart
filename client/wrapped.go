package client

import (
	"github.com/gocaine/go-dart/hardware"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type WrappedClient struct {
	CurrentGameId   int
	LatestGameState GameState
	client          *DartClient
}

func NewWrappedClient() *WrappedClient {
	return &WrappedClient{client: NewClient()}
}

func (wrapped *WrappedClient) Consume(event hardware.InputEvent) {
	// simply delegate to the client
	board := viper.GetString("board")
	log.Printf("Board: %s", board)
	wrapped.client.FireDart(board, event.Sector, event.Multiplier)
}

func (wrapped *WrappedClient) Shutdown() {
	log.Info("Mocked client is shuting down")
}
