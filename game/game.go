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
	// BoardHasLeft is call to notify the game a board has been disconnected. Returns true if the game continues despite this event
	BoardHasLeft(board string) bool
	// HoldOrNextPlayer switch game state between ONHOLD and PLAYING with side effects according to game implementation
	HoldOrNextPlayer()
	nextPlayer()
	nextDart()
}

// BaseGame common Game struct
type BaseGame struct {
	state        *common.GameState
	DisplayStyle string
	rank         int
}

// State : get the current GameState
func (game *BaseGame) State() *common.GameState {
	return game.state
}

// Start start the game, Darts will be handled
func commonStart(game Game) (error error) {
	if game.State().Ongoing == common.READY && len(game.State().Players) > 0 {
		state := game.State()
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

// BoardHasLeft is call to notify the game a board has been disconnected. It returns true if the game continues despite this event.
func (game *BaseGame) BoardHasLeft(board string) bool {
	for _, p := range game.state.Players {
		if p.Board == board {
			log.Infof("game is over because the board %s from player %s has been disconnected", board, p.Name)
			// end the game has one player has left
			game.state.Ongoing = common.OVER
			game.state.LastMsg = "Board " + board + " has been disconnected"
			return false
		}
	}
	return true
}

// AddPlayer add a new player to the game
func commonAddPlayer(game Game, board string, name string) (error error) {
	if game.State().Ongoing == common.INITIALIZING || game.State().Ongoing == common.READY {
		for _, p := range game.State().Players {
			if name == p.Name {
				// player with same name is already registred
				return errors.New("Player name is already in use")
			}
		}

		log.WithFields(log.Fields{"player": name, "board": board}).Infof("Player added to the game")

		game.State().Players = append(game.State().Players, common.PlayerState{Name: name, Board: board, Visits: make([]common.Sector, 0, 3)})

		// now that we have at least one player, we are in a ready state, waiting for other players or the first dart
		game.State().Ongoing = common.READY
	} else {
		error = errors.New("Player cannot be started")
	}
	return
}

func commonNextDart(game Game) {
	state := game.State()
	if state.CurrentDart == 2 {
		game.HoldOrNextPlayer()
	} else {
		state.CurrentDart++
		log.WithFields(log.Fields{"player": state.CurrentPlayer, "dart": state.CurrentDart}).Info("One more dart")
	}
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING with side effects according to game implementation
func commonHoldOrNextPlayer(game Game) {
	if game.State().Ongoing == common.PLAYING || game.State().Ongoing == common.READY {
		game.State().Ongoing = common.ONHOLD
		game.State().LastMsg = "Next Player"
		game.State().LastSector = common.Sector{}
	} else if game.State().Ongoing == common.ONHOLD {
		game.State().Ongoing = common.PLAYING
		game.State().LastMsg = ""
		game.nextPlayer()
	}
}

func commonNextPlayer(game Game) {
	state := game.State()

	// reset visits
	state.Players[state.CurrentPlayer].Visits = make([]common.Sector, 0, 3)

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

func commonHandleDartChecks(game Game, sector common.Sector) (error error) {

	if game.State().Ongoing == common.ONHOLD {
		error = errors.New("Game is on hold and not ready to handle darts")
		return
	}

	if game.State().Ongoing == common.READY {
		// first dart starts the game
		error = game.Start()
		if error != nil {
			return
		}
	}

	if game.State().Ongoing != common.PLAYING {
		error = errors.New("Game is not started or is ended")
		return
	}

	if !sector.IsValid() {
		log.WithFields(log.Fields{"sector": sector}).Error("Invalid sector")
		error = errors.New("Sector is not a valid one")
		return
	}

	return
}
