package server

import (
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
}

func (game *AGame) GetState() *common.GameState {
	return game.State
}
