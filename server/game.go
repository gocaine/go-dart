package server

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"go-dart/common"
)

type Game interface {
	Start() error
	AddPlayer(name string) error
	HandleDart(sector common.Sector) (*common.GameState, error)
	GetState() *common.GameState
}

type AGame struct {
	State *common.GameState
	rank  int
}

func (game *AGame) GetState() *common.GameState {
	return game.State
}

func (game *AGame) Start() (error error) {
	if game.State.Ongoing == common.READY && len(game.State.Players) > 0 {
		state := game.State
		state.Ongoing = common.PLAYING
		state.CurrentPlayer = 0
		state.CurrentDart = 0
		for i := range state.Players {
			state.Players[i].Score = 0
		}
		state.Round = 1
		log.Infof("The game is now started")
	} else {
		error = errors.New("Game cannot start")
	}
	return
}
