package game

import (
	"errors"
	"fmt"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
)

// Highest is a highscore Game, within a fixed number of visit, winner is the highest score
type Highest struct {
	AGame
	rounds int
}

// OptionHighest is the struct to handle GameHighest parameters
type OptionHighest struct {
	Rounds int
}

// NewGameHighest : GameHighest constructor using a OptionHighest
func NewGameHighest(board string, opt OptionHighest) *Highest {

	g := new(Highest)
	g.SetBoard(board)
	g.rounds = opt.Rounds
	g.state = common.NewGameState()

	g.DisplayStyle = fmt.Sprintf("%d visits HighScore", opt.Rounds)

	return g
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
func (game *Highest) HandleDart(sector common.Sector) (result *common.GameState, error error) {

	if game.state.Ongoing == common.READY {
		// first dart starts the game
		err := game.Start()
		if err != nil {
			error = err
			return
		}
	}

	if game.state.Ongoing != common.PLAYING {
		error = errors.New("Game is not started or is ended")
		return
	}

	if !sector.IsValid() {
		log.WithFields(log.Fields{"sector": sector}).Error("Invalid sector")
		error = errors.New("Sector is not a valid one")
		return
	}

	point := sector.Val * sector.Pos
	state := game.state

	state.LastSector = sector

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "score": point}).Info("Scored")

	state.Players[state.CurrentPlayer].Score += point

	log.WithFields(log.Fields{"state.Round": state.Round, "game.rounds": game.rounds}).Info("Rounds")
	if state.Round == game.rounds && state.CurrentDart == 2 {
		game.winner()
		if game.state.Ongoing == common.PLAYING {
			game.nextPlayer()
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
