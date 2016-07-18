package client

import (
	"fmt"
	"github.com/gocaine/go-dart/hardware"

	log "github.com/Sirupsen/logrus"
)

type WrappedClient struct {
	CurrentGameId   int
	LatestGameState GameState
	client          *DartClient
}

func NewWrappedClient() *WrappedClient {
	return &WrappedClient{client: NewClient()}
}

func (wrapped *WrappedClient) Prepare() error {
	fmt.Print("Select game style: ")
	var gameStyle string
	fmt.Scanln(&gameStyle)
	game, err := wrapped.client.CreateGame(gameStyle)
	if err != nil {
		return err
	}
	wrapped.CurrentGameId = game.Id
	wrapped.LatestGameState = game.Game

	// Simulate 2 players
	_, err = wrapped.client.CreatePlayer(wrapped.CurrentGameId, "player 1")
	if err != nil {
		return err
	}
	_, err = wrapped.client.CreatePlayer(wrapped.CurrentGameId, "player 2")
	if err != nil {
		return err
	}
	log.Info("Ready to dart !")

	return nil
}

func (wrapped *WrappedClient) Consume(event hardware.InputEvent) {
	// simply delegate to the client
	wrapped.client.FireDart(wrapped.CurrentGameId, event.Sector, event.Multiplier)
}

func (wrapped *WrappedClient) Shutdown() {
	log.Info("Mocked client is shuting down")
}
