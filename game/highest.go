package game

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/i18n"
)

// Highest is a highscore Game, within a fixed number of visit, winner is the highest score
type Highest struct {
	BaseGame
	rounds int
}

// OptionHighest is the struct to handle GameHighest parameters
type OptionHighest struct {
	Rounds int
}

// NewGameHighest : GameHighest constructor using a OptionHighest
func NewGameHighest(ctx common.GameContext, opts map[string]interface{}) (g *Highest, err error) {
	opt := newOptionHighest(opts)
	if opt.Rounds < 1 {
		err = errors.New(i18n.Translation("game.highest.error.rounds", ctx.Locale))
		return
	}
	g = new(Highest)
	g.rounds = opt.Rounds
	g.state = common.NewGameState()

	g.DisplayStyle = fmt.Sprintf(i18n.Translation("game.highest.display", ctx.Locale), opt.Rounds)

	return
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
func (game *Highest) HandleDart(ctx common.GameContext, sector common.Sector) (result *common.GameState, error error) {

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

	log.WithFields(log.Fields{"state.Round": state.Round, "game.rounds": game.rounds}).Info("Rounds")
	if state.Round == game.rounds && state.CurrentDart == 2 {
		game.winner()
		if game.state.Ongoing == common.PLAYING {
			game.HoldOrNextPlayer(ctx)
		}

	} else {
		game.NextDart(ctx)
	}
	result = state
	return
}

func (game *Highest) winner() {
	state := game.state
	if game.state.CurrentPlayer == len(state.Players)-1 {
		game.state.Ongoing = common.OVER
		sort.Sort(common.ByScore(state.Players))
		for i := 0; i < len(state.Players); i++ {
			state.Players[i].Rank = i + 1
		}
	}
}

// Start start the game, Darts will be handled
func (game *Highest) Start(ctx common.GameContext) error {
	return commonStart(ctx, game)
}

// AddPlayer add a new player to the game
func (game *Highest) AddPlayer(ctx common.GameContext, board string, name string) error {
	return commonAddPlayer(ctx, game, board, name)
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *Highest) HoldOrNextPlayer(ctx common.GameContext) {
	commonHoldOrNextPlayer(ctx, game)
}

// NextDart is called after each dart when the same palyer play again
func (game *Highest) NextDart(ctx common.GameContext) {
	commonNextDart(ctx, game)
}

// NextPlayer is called when the current player end his visit
func (game *Highest) NextPlayer(ctx common.GameContext) {
	commonNextPlayer(ctx, game)
}

var gsHighestOptions = []common.GameOption{{"Rounds", "int", "game.highest.options.rounds", 5}}

// GsHighest GameStyle for Highest series
var GsHighest = common.GameStyle{
	"game.highest.name",
	"HIGHEST",
	"game.highest.rules",
	gsHighestOptions}

func newOptionHighest(opts map[string]interface{}) OptionHighest {
	o := OptionHighest{}
	gameOptionFiller(&o, gsHighestOptions, opts)
	return o
}
