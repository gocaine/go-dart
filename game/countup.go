package game

import (
	"errors"
	"sort"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/i18n"
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
func NewGameCountUp(ctx common.GameContext, opts map[string]interface{}) (g *CountUp, err error) {
	opt := newOptionCountUp(opts)
	if opt.Target < 61 {
		err = errors.New(i18n.Translation("game.countup.error.target", ctx.Locale))
		return
	}
	g = new(CountUp)
	g.target = opt.Target
	g.state = common.NewGameState()

	g.DisplayStyle = fmt.Sprintf(i18n.Translation("game.countup.display", ctx.Locale), opt.Target)

	return
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
func (game *CountUp) HandleDart(ctx common.GameContext, sector common.Sector) (result *common.GameState, error error) {

	error = commonHandleDartChecks(ctx, game, sector)
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
		game.winner(ctx)
		if game.state.Ongoing == common.PLAYING {
			game.HoldOrNextPlayer(ctx)
		}

	} else {
		game.NextDart(ctx)
	}
	result = state
	return
}

func (game *CountUp) winner(ctx common.GameContext) {
	state := game.state
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	ctx.MessageHandler("game.message.rank", state.CurrentPlayer, game.rank+1)
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
func (game *CountUp) AddPlayer(ctx common.GameContext, board string, name string) error {
	return commonAddPlayer(ctx, game, board, name)
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *CountUp) HoldOrNextPlayer(ctx common.GameContext) {
	commonHoldOrNextPlayer(ctx, game)
}

// Start start the game, Darts will be handled
func (game *CountUp) Start(ctx common.GameContext) error {
	return commonStart(ctx, game)
}

// NextDart is called after each dart when the same palyer play again
func (game *CountUp) NextDart(ctx common.GameContext) {
	commonNextDart(ctx, game)
}

// NextPlayer is called when the current player end his visit
func (game *CountUp) NextPlayer(ctx common.GameContext) {
	commonNextPlayer(ctx, game)
}

var gsCountUpOptions = []common.GameOption{{"Target", "int", "game.countup.options.target", 500}}

// GsCountUp GameStyle for CountUp series
var GsCountUp = common.GameStyle{
	"game.countup.name",
	"COUNTUP",
	"game.countup.rules",
	gsCountUpOptions}

func newOptionCountUp(opts map[string]interface{}) OptionCountUp {
	o := OptionCountUp{}
	gameOptionFiller(&o, gsCountUpOptions, opts)
	return o
}
