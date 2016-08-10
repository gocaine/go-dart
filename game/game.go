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
	AddPlayer(name string) error
	// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
	HandleDart(sector common.Sector) (*common.GameState, error)
	// GetState, get the current GameState
	GetState() *common.GameState
	Board() string
	SetBoard(board string)
}

// AGame common Game struct
type AGame struct {
	State        *common.GameState
	DisplayStyle string
	rank         int
	board        string
}

// GetState : get the current GameState
func (game *AGame) GetState() *common.GameState {
	return game.State
}

// Start start the game, Darts will be handled
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

// AddPlayer add a new player to the game
func (game *AGame) AddPlayer(name string) (error error) {
	if game.State.Ongoing == common.INITIALIZING || game.State.Ongoing == common.READY {
		log.WithFields(log.Fields{"player": name}).Infof("Player added to the game")
		game.State.Players = append(game.State.Players, common.PlayerState{Name: name})
		// now that we have at least one player, we are in a ready state, waiting for other players or the first dart
		game.State.Ongoing = common.READY
	} else {
		error = errors.New("Game cannot be started")
	}
	return
}

func (game *AGame) nextDart() {
	state := game.State
	if state.CurrentDart == 2 {
		game.nextPlayer()
	} else {
		state.CurrentDart++
		log.WithFields(log.Fields{"player": state.CurrentPlayer, "dart": state.CurrentDart}).Info("One more dart")
	}
}

func (game *AGame) nextPlayer() {
	state := game.State
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

// Board is the board getter
func (game *AGame) Board() string {
	return game.board
}

// SetBoard is the board setter
func (game *AGame) SetBoard(board string) {
	game.board = board
}
