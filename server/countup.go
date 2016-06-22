package server

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"go-dart/common"
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

	return g
}

func (game *GameCountUp) AddPlayer(name string) (error error) {
	if game.State.Ongoing == common.INITIALIZING || game.State.Ongoing == common.READY {
		log.WithFields(log.Fields{"player": name}).Infof("Player added to the game")
		game.State.Players = append(game.State.Players, common.PlayerState{Name: name})
		// now that we have at least one player, we are in a ready state, waiting for other players or the first dart
		game.State.Ongoing = common.READY
	} else {
		error = errors.New("Game cannot be started")
	}
	return
}

func (game *GameCountUp) Start() (error error) {
	if game.State.Ongoing == common.READY && len(game.State.Players) > 0 {
		state := game.State
		state.Ongoing = common.PLAYING
		state.CurrentPlayer = 0
		state.CurrentDart = 0
		for i := range state.Players {
			state.Players[i].Score = 0
		}
		state.Round = 1
		log.Infof("The game is now started")
	} else {
		error = errors.New("Game cannot start")
	}
	return
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

func (game *GameCountUp) nextPlayer() {
	state := game.State
	state.CurrentDart = 0
	state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Players)
	for state.Players[state.CurrentPlayer].Score >= game.target {
		state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Players)
	}
	log.WithFields(log.Fields{"player": state.CurrentPlayer}).Info("Next player")
}
