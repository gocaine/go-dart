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
	AGame
	target int
}

// OptionCountUp is the struct to handle GameCountUp parameters
type OptionCountUp struct {
	Target int
}

// NewGameCountUp : GameCountUp constructor using a OptionCountUp
func NewGameCountUp(board string, opt OptionCountUp) *CountUp {

	g := new(CountUp)
	g.SetBoard(board)
	g.target = opt.Target
	g.state = common.NewGameState()

	g.DisplayStyle = fmt.Sprintf("Count-Up %d", opt.Target)

	return g
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
func (game *CountUp) HandleDart(sector common.Sector) (result *common.GameState, error error) {

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

	if state.Players[state.CurrentPlayer].Score >= game.target {
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
