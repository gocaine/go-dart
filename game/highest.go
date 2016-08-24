package game

import (
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
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
func NewGameHighest(opt OptionHighest) *Highest {

	g := new(Highest)
	g.rounds = opt.Rounds
	g.state = common.NewGameState()

	g.DisplayStyle = fmt.Sprintf("%d visits HighScore", opt.Rounds)

	return g
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
func (game *Highest) HandleDart(sector common.Sector) (result *common.GameState, error error) {

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

	log.WithFields(log.Fields{"state.Round": state.Round, "game.rounds": game.rounds}).Info("Rounds")
	if state.Round == game.rounds && state.CurrentDart == 2 {
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
func (game *Highest) Start() error {
	return commonStart(game)
}

// AddPlayer add a new player to the game
func (game *Highest) AddPlayer(board string, name string) error {
	return commonAddPlayer(game, board, name)
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *Highest) HoldOrNextPlayer() {
	commonHoldOrNextPlayer(game)
}

func (game *Highest) nextDart() {
	commonNextDart(game)
}

func (game *Highest) nextPlayer() {
	commonNextPlayer(game)
}
