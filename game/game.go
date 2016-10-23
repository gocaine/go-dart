package game

import (
	"errors"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/i18n"
	"reflect"
)

// Game interface, should be implemented by all game (rules) implems
type Game interface {
	// Start start the game, Darts will be handled
	Start(ctx common.GameContext) error
	// AddPlayer add a new player to the game
	AddPlayer(ctx common.GameContext, board string, name string) error
	// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
	HandleDart(ctx common.GameContext, sector common.Sector) (*common.GameState, error)
	// GetState, get the current GameState
	State() *common.GameState
	// BoardHasLeft is call to notify the game a board has been disconnected. Returns true if the game continues despite this event
	BoardHasLeft(ctx common.GameContext, board string) bool
	// HoldOrNextPlayer switch game state between ONHOLD and PLAYING with side effects according to game implementation
	HoldOrNextPlayer(ctx common.GameContext)
	nextPlayer(ctx common.GameContext)
	nextDart(ctx common.GameContext)
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
func commonStart(ctx common.GameContext, game Game) (error error) {
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
		error = errors.New(i18n.Translation("game.error.cantstart", ctx.Locale))
	}
	return
}

// BoardHasLeft is call to notify the game a board has been disconnected. It returns true if the game continues despite this event.
func (game *BaseGame) BoardHasLeft(ctx common.GameContext, board string) bool {
	for _, p := range game.state.Players {
		if p.Board == board {
			log.Infof("game is over because the board %s from player %s has been disconnected", board, p.Name)
			// end the game has one player has left
			game.state.Ongoing = common.OVER
			game.state.LastMsg = i18n.Translation("game.message.disconnect", ctx.Locale)
			return false
		}
	}
	return true
}

// AddPlayer add a new player to the game
func commonAddPlayer(ctx common.GameContext, game Game, board string, name string) (error error) {
	if game.State().Ongoing == common.INITIALIZING || game.State().Ongoing == common.READY {
		for _, p := range game.State().Players {
			if name == p.Name {
				// player with same name is already registred
				return errors.New(i18n.Translation("game.message.player.exists", ctx.Locale))
			}
		}

		log.WithFields(log.Fields{"player": name, "board": board}).Infof("Player added to the game")

		game.State().Players = append(game.State().Players, common.PlayerState{Name: name, Board: board, Visits: make([]common.Sector, 0, 3)})

		// now that we have at least one player, we are in a ready state, waiting for other players or the first dart
		game.State().Ongoing = common.READY
	} else {
		error = errors.New(i18n.Translation("game.message.player.notadded", ctx.Locale))
	}
	return
}

func commonNextDart(ctx common.GameContext, game Game) {
	state := game.State()
	if state.CurrentDart == 2 {
		game.HoldOrNextPlayer(ctx)
	} else {
		state.CurrentDart++
		log.WithFields(log.Fields{"player": state.CurrentPlayer, "dart": state.CurrentDart}).Info("One more dart")
	}
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING with side effects according to game implementation
func commonHoldOrNextPlayer(ctx common.GameContext, game Game) {
	if game.State().Ongoing == common.PLAYING || game.State().Ongoing == common.READY {
		game.State().Ongoing = common.ONHOLD
		game.State().LastMsg = i18n.Translation("game.message.player.next", ctx.Locale)
		game.State().LastSector = common.Sector{}
	} else if game.State().Ongoing == common.ONHOLD {
		game.State().Ongoing = common.PLAYING
		game.State().LastMsg = ""
		game.nextPlayer(ctx)
	}
}

func commonNextPlayer(ctx common.GameContext, game Game) {
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

func commonHandleDartChecks(ctx common.GameContext, game Game, sector common.Sector) (error error) {

	if game.State().Ongoing == common.ONHOLD {
		error = errors.New(i18n.Translation("game.error.onhold", ctx.Locale))
		return
	}

	if game.State().Ongoing == common.READY {
		// first dart starts the game
		error = game.Start(ctx)
		if error != nil {
			return
		}
	}

	if game.State().Ongoing != common.PLAYING {
		error = errors.New(i18n.Translation("game.error.notstarted", ctx.Locale))
		return
	}

	if !sector.IsValid() {
		log.WithFields(log.Fields{"sector": sector}).Error("Invalid sector")
		error = errors.New(i18n.Translation("game.error.sector.invalid", ctx.Locale))
		return
	}

	return
}

func gameOptionFiller(o interface{}, gsOpts []common.GameOption, opts map[string]interface{}) error {
	v := reflect.ValueOf(o).Elem()
	for _, val := range gsOpts {
		f := v.FieldByName(val.Name)
		f.Set(reflect.ValueOf(val.Default))
		if iv, ok := opts[val.Name]; ok {
			if ivv := reflect.ValueOf(iv); ivv.Type().ConvertibleTo(f.Type()) {
				f.Set(ivv.Convert(f.Type()))
			} else {
				return fmt.Errorf("%v is an invalid value for %s", iv, val.Name)
			}
		}
	}
	return nil
}

// Flavors gives translated game styles and flavors (rules and options...)
func Flavors(ctx common.GameContext) []common.GameStyle {
	games := []common.GameStyle{GsX01, GsCountUp, GsHighest, GsCricket}
	for idx1, gs := range games {
		games[idx1].Name = i18n.Translation(gs.Name, ctx.Locale)
		games[idx1].Rules = i18n.Translation(gs.Rules, ctx.Locale)
		for idx2, opt := range gs.Options {
			games[idx1].Options[idx2].Desc = i18n.Translation(opt.Desc, ctx.Locale)
		}
	}
	return games
}
