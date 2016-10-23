package game

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/i18n"
)

// Gamex01 is a x01 series Game (301, 501-Double-Out, ...)
type Gamex01 struct {
	BaseGame
	score     int
	doubleOut bool
	accu      int
}

// Optionx01 is the struct to handle Gamex01 parameters
type Optionx01 struct {
	Score     int
	DoubleOut bool
}

// NewGamex01 : Gamex01 constructor
func NewGamex01(ctx common.GameContext, opts map[string]interface{}) (g *Gamex01, err error) {
	opt, err := newOptionx01(opts)
	if err != nil {
		return
	}
	if opt.Score < 61 {
		err = errors.New(i18n.Translation("game.x01.error.score", ctx.Locale))
		return
	}
	g = new(Gamex01)
	g.doubleOut = opt.DoubleOut
	g.score = opt.Score
	g.state = common.NewGameState()

	dStyle := "game.x01.display.x01"
	if opt.DoubleOut {
		dStyle = "game.x01.display.doubleout"
	}
	g.DisplayStyle = fmt.Sprintf(i18n.Translation(dStyle, ctx.Locale), opt.Score)

	return
}

// Start start the game, Darts will be handled
func (game *Gamex01) Start(ctx common.GameContext) (error error) {
	if game.state.Ongoing == common.READY && len(game.state.Players) > 0 && game.score > 0 {
		state := game.state
		state.Ongoing = common.PLAYING
		state.CurrentPlayer = 0
		state.CurrentDart = 0
		for i := range state.Players {
			state.Players[i].Score = game.score
		}
		log.Infof("The game is now started")
	} else {
		error = errors.New(i18n.Translation("game.error.cantstart", ctx.Locale))
	}
	return
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules of x01, and the context. Return a GameState
func (game *Gamex01) HandleDart(ctx common.GameContext, sector common.Sector) (result *common.GameState, error error) {

	error = commonHandleDartChecks(ctx, game, sector)
	if error != nil {
		return
	}

	point := sector.Val * sector.Pos
	game.accu += point
	state := game.state

	state.LastSector = sector

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "score": point}).Info("Scored")

	state.Players[state.CurrentPlayer].Score -= point

	log.Info("Current dart ", game.state.CurrentDart)
	state.Players[state.CurrentPlayer].Visits = append(state.Players[state.CurrentPlayer].Visits, sector)

	if state.Players[state.CurrentPlayer].Score > 0 {
		if game.doubleOut && state.Players[state.CurrentPlayer].Score == 1 {
			state.LastMsg = i18n.Translation("game.x01.message.doubleout", ctx.Locale)
			game.resetVisit()
			game.HoldOrNextPlayer(ctx)
		} else {
			game.nextDart(ctx)
		}

	} else if state.Players[state.CurrentPlayer].Score == 0 {
		if game.doubleOut && sector.Pos != 2 {
			state.LastMsg = i18n.Translation("game.x01.message.doubleout", ctx.Locale)
			game.resetVisit()
			game.HoldOrNextPlayer(ctx)
		} else {
			game.winner(ctx)
			if game.state.Ongoing == common.PLAYING {
				game.HoldOrNextPlayer(ctx)
			}
		}

	} else {
		state.LastMsg = i18n.Translation("game.x01.message.overscore", ctx.Locale)
		game.resetVisit()
		game.HoldOrNextPlayer(ctx)
	}
	result = state
	return
}

func (game *Gamex01) winner(ctx common.GameContext) {
	state := game.state
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	state.LastMsg = fmt.Sprintf(i18n.Translation("game.message.rank", ctx.Locale), state.CurrentPlayer, game.rank+1)
	game.rank++
	if game.rank >= len(state.Players)-1 {
		game.state.Ongoing = common.OVER
		sort.Sort(common.ByRank(state.Players))
		if len(state.Players) > 1 {
			state.Players[len(state.Players)-1].Rank = game.rank + 1
		}
	}
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *Gamex01) HoldOrNextPlayer(ctx common.GameContext) {
	commonHoldOrNextPlayer(ctx, game)
}

// AddPlayer add a new player to the game
func (game *Gamex01) AddPlayer(ctx common.GameContext, board string, name string) error {
	return commonAddPlayer(ctx, game, board, name)
}

func (game *Gamex01) nextPlayer(ctx common.GameContext) {

	game.accu = 0
	state := game.State()

	// reset visits
	state.Players[state.CurrentPlayer].Visits = make([]common.Sector, 0, 3)

	state.CurrentDart = 0
	state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Players)
	for state.Players[state.CurrentPlayer].Score == 0 {
		state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Players)
	}
	log.WithFields(log.Fields{"player": state.CurrentPlayer}).Info("Next player")
}

func (game *Gamex01) nextDart(ctx common.GameContext) {
	state := game.state
	if state.CurrentDart == 2 {
		game.HoldOrNextPlayer(ctx)
	} else {
		state.CurrentDart++
		log.WithFields(log.Fields{"player": state.CurrentPlayer, "dart": state.CurrentDart}).Info("One more dart")
	}
}

func (game *Gamex01) resetVisit() {
	state := game.state
	state.Players[state.CurrentPlayer].Score += game.accu
}

var gsX01Options = []common.GameOption{
	{"Score", "int", "game.x01.options.score", 501},
	{"DoubleOut", "bool", "game.x01.options.doubleout", false}}

// GsX01 GameStyle for X01 series
var GsX01 = common.GameStyle{
	"game.x01.name",
	"X01",
	"game.x01.rules",
	gsX01Options}

func newOptionx01(opts map[string]interface{}) (o Optionx01, err error) {
	o = Optionx01{}
	err = gameOptionFiller(&o, gsX01Options, opts)

	return
}
