package game

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
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
func NewGamex01(opts map[string]interface{}) (g *Gamex01, err error) {
	opt := newOptionx01(opts)
	if opt.Score < 61 {
		err = errors.New("Score should be at least 61")
	}
	g = new(Gamex01)
	g.doubleOut = opt.DoubleOut
	g.score = opt.Score
	g.state = common.NewGameState()

	dStyle := ""
	if opt.DoubleOut {
		dStyle = " Double-Out"
	}
	g.DisplayStyle = fmt.Sprintf("%d%s", opt.Score, dStyle)

	return
}

// Start start the game, Darts will be handled
func (game *Gamex01) Start() (error error) {
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
		error = errors.New("Game cannot start")
	}
	return
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules of x01, and the context. Return a GameState
func (game *Gamex01) HandleDart(sector common.Sector) (result *common.GameState, error error) {

	error = commonHandleDartChecks(game, sector)
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
			state.LastMsg = "You should end with a double"
			game.resetVisit()
			game.HoldOrNextPlayer()
		} else {
			game.nextDart()
		}

	} else if state.Players[state.CurrentPlayer].Score == 0 {
		if game.doubleOut && sector.Pos != 2 {
			state.LastMsg = "You should end with a double"
			game.resetVisit()
			game.HoldOrNextPlayer()
		} else {
			game.winner()
			if game.state.Ongoing == common.PLAYING {
				game.HoldOrNextPlayer()
			}
		}

	} else {
		state.LastMsg = "You went beyond the target dude !"
		game.resetVisit()
		game.HoldOrNextPlayer()
	}
	result = state
	return
}

func (game *Gamex01) winner() {
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

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *Gamex01) HoldOrNextPlayer() {
	commonHoldOrNextPlayer(game)
}

// AddPlayer add a new player to the game
func (game *Gamex01) AddPlayer(board string, name string) error {
	return commonAddPlayer(game, board, name)
}

func (game *Gamex01) nextPlayer() {

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

func (game *Gamex01) nextDart() {
	state := game.state
	if state.CurrentDart == 2 {
		game.HoldOrNextPlayer()
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
	{"Score", "int", "The score from which to reach 0", 501},
	{"DoubleOut", "bool", "If set to true, the players have to end with a double (and so reach 0)", false}}

// GsX01 GameStyle for X01 series
var GsX01 = common.GameStyle{
	"X01 : 301, 501,...",
	"X01",
	"All players start with the same points (301 / 501 / ...) and attempt to reach zero. " +
		"If a player scores more than the total required to reach zero, " +
		"the player \"busts\" and the score returns to the score that was existing at the start of the turn.",
	gsX01Options}

func newOptionx01(opts map[string]interface{}) Optionx01 {
	o := Optionx01{}
	gameOptionFiller(&o, gsX01Options, opts)

	return o
}
