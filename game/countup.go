package game

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
)

// CountUp is a highscore Game, winner is the first to obtain a given score or more
type CountUp struct {
	BaseGame
	target int
}

// OptionCountUp is the struct to handle GameCountUp parameters
type OptionCountUp struct {
	Target int
}

// NewGameCountUp : GameCountUp constructor using a OptionCountUp
func NewGameCountUp(opts map[string]interface{}) (g *CountUp, err error) {
	opt := newOptionCountUp(opts)
	if opt.Target < 61 {
		err = errors.New("Target should be at least 61")
		return
	}
	g = new(CountUp)
	g.target = opt.Target
	g.state = common.NewGameState()

	g.DisplayStyle = fmt.Sprintf("Count-Up %d", opt.Target)

	return
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
func (game *CountUp) HandleDart(sector common.Sector) (result *common.GameState, error error) {

	error = commonHandleDartChecks(game, sector)
	if error != nil {
		return
	}

	point := sector.Val * sector.Pos
	state := game.state

	state.LastSector = sector

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "score": point}).Info("Scored")

	state.Players[state.CurrentPlayer].Score += point

	state.Players[state.CurrentPlayer].Visits = append(state.Players[state.CurrentPlayer].Visits, sector)

	if state.Players[state.CurrentPlayer].Score >= game.target {
		game.winner()
		if game.state.Ongoing == common.PLAYING {
			game.HoldOrNextPlayer()
		}

	} else {
		game.nextDart()
	}
	result = state
	return
}

func (game *CountUp) winner() {
	state := game.state
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	state.LastMsg = fmt.Sprintf("Player %d end at rank #%d", state.CurrentPlayer, game.rank+1)
	game.rank++
	if game.rank >= len(state.Players)-1 {
		game.state.Ongoing = common.OVER
		sort.Sort(common.ByRank(state.Players))
		if len(state.Players) > 1 {
			state.Players[len(state.Players)-1].Rank = game.rank + 1
		}
	}
}

// AddPlayer add a new player to the game
func (game *CountUp) AddPlayer(board string, name string) error {
	return commonAddPlayer(game, board, name)
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *CountUp) HoldOrNextPlayer() {
	commonHoldOrNextPlayer(game)
}

// Start start the game, Darts will be handled
func (game *CountUp) Start() error {
	return commonStart(game)
}

func (game *CountUp) nextDart() {
	commonNextDart(game)
}

func (game *CountUp) nextPlayer() {
	commonNextPlayer(game)
}

var gsCountUpOptions = []common.GameOption{{"Target", "int", "The score to reach", 500}}

// GsCountUp GameStyle for CountUp series
var GsCountUp = common.GameStyle{
	"CountUp",
	"COUNTUP",
	"All players start with 0 points and attempt to reach the given target (300 / 500 / ...). ",
	gsCountUpOptions}

func newOptionCountUp(opts map[string]interface{}) OptionCountUp {
	o := OptionCountUp{}
	gameOptionFiller(&o, gsCountUpOptions, opts)
	return o
}
