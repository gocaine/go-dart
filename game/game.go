package game

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
)

// Game interface, should be implemented by all game (rules) implems
type Game interface {
	// Start start the game, Darts will be handled
	Start() error
	// AddPlayer add a new player to the game
	AddPlayer(board string, name string) error
	// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
	HandleDart(sector common.Sector) (*common.GameState, error)
	// GetState, get the current GameState
	State() *common.GameState
}

// AGame common Game struct
type AGame struct {
	state        *common.GameState
	DisplayStyle string
	rank         int
}

// State : get the current GameState
func (game *AGame) State() *common.GameState {
	return game.state
}

// Start start the game, Darts will be handled
func (game *AGame) Start() (error error) {
	if game.state.Ongoing == common.READY && len(game.state.Players) > 0 {
		state := game.state
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

// AddPlayer add a new player to the game
func (game *AGame) AddPlayer(board string, name string) (error error) {
	if game.state.Ongoing == common.INITIALIZING || game.state.Ongoing == common.READY {
		log.WithFields(log.Fields{"player": name, "board": board}).Infof("Player added to the game")
		game.state.Players = append(game.state.Players, common.PlayerState{Name: name, Board: board})
		// now that we have at least one player, we are in a ready state, waiting for other players or the first dart
		game.state.Ongoing = common.READY
	} else {
		error = errors.New("Game cannot be started")
	}
	return
}

func (game *AGame) nextDart() {
	state := game.state
	if state.CurrentDart == 2 {
		game.nextPlayer()
	} else {
		state.CurrentDart++
		log.WithFields(log.Fields{"player": state.CurrentPlayer, "dart": state.CurrentDart}).Info("One more dart")
	}
}

func (game *AGame) nextPlayer() {
	state := game.state
	state.CurrentDart = 0
	state.CurrentPlayer = state.CurrentPlayer + 1
	if state.CurrentPlayer >= len(state.Players) {
		state.CurrentPlayer = 0
		state.Round++
	}
	for state.Players[state.CurrentPlayer].Rank > 0 {
		state.CurrentPlayer = state.CurrentPlayer + 1
		if state.CurrentPlayer >= len(state.Players) {
			state.CurrentPlayer = 0
			state.Round++
		}
	}
	log.WithFields(log.Fields{"player": state.CurrentPlayer}).Info("Next player")
}
