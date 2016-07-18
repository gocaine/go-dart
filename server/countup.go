package server

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"sort"
)

type GameCountUp struct {
	AGame
	target int
}

type OptionCountUp struct {
	Target int
}

func NewGameCountUp(opt OptionCountUp) *GameCountUp {

	g := new(GameCountUp)
	g.target = opt.Target
	g.State = common.NewGameState()

	g.DisplayStyle = fmt.Sprintf("Count-Up %d", opt.Target)

	return g
}

func (game *GameCountUp) HandleDart(sector common.Sector) (result *common.GameState, error error) {

	if game.State.Ongoing == common.READY {
		// first dart starts the game
		err := game.Start()
		if err != nil {
			error = err
			return
		}
	}

	if game.State.Ongoing != common.PLAYING {
		error = errors.New("Game is not started or is ended")
		return
	}

	if !sector.IsValid() {
		log.WithFields(log.Fields{"sector": sector}).Error("Invalid sector")
		error = errors.New("Sector is not a valid one")
		return
	}

	point := sector.Val * sector.Pos
	state := game.State

	state.LastSector = sector

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "score": point}).Info("Scored")

	state.Players[state.CurrentPlayer].Score += point

	if state.Players[state.CurrentPlayer].Score >= game.target {
		game.winner()
		if game.State.Ongoing == common.PLAYING {
			game.nextPlayer()
		}

	} else {
		game.nextDart()
	}
	result = state
	return
}

func (game *GameCountUp) winner() {
	state := game.State
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	state.LastMsg = fmt.Sprintf("Player %d end at rank #%d", state.CurrentPlayer, game.rank+1)
	game.rank++
	if game.rank >= len(state.Players)-1 {
		game.State.Ongoing = common.OVER
		sort.Sort(common.ByRank(state.Players))
		if len(state.Players) > 1 {
			state.Players[len(state.Players)-1].Rank = game.rank + 1
		}
	}
}
