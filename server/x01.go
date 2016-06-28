package server

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"go-dart/common"
	"sort"
)

type Gamex01 struct {
	AGame
	score     int
	doubleOut bool
	accu      int
}

type Optionx01 struct {
	Score     int
	DoubleOut bool
}

func NewGamex01(opt Optionx01) *Gamex01 {
	g := new(Gamex01)
	g.doubleOut = opt.DoubleOut
	g.score = opt.Score
	g.State = common.NewGameState()

	dStyle := ""
	if opt.DoubleOut {
		dStyle = " Double-Out"
	}
	g.DisplayStyle = fmt.Sprintf("%d%s", opt.Score, dStyle)

	return g
}

func (game *Gamex01) AddPlayer(name string) (error error) {
	if game.State.Ongoing == common.INITIALIZING || game.State.Ongoing == common.READY {
		log.WithFields(log.Fields{"player": name}).Infof("Player added to the game")
		game.State.Players = append(game.State.Players, common.PlayerState{Name: name, Score: game.score})
		// now that we have at least one player, we are in a ready state, waiting for other players or the first dart
		game.State.Ongoing = common.READY
	} else {
		error = errors.New("Game cannot be started")
	}
	return
}

func (game *Gamex01) Start() (error error) {
	if game.State.Ongoing == common.READY && len(game.State.Players) > 0 && game.score > 0 {
		state := game.State
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

func (game *Gamex01) HandleDart(sector common.Sector) (result *common.GameState, error error) {

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
	game.accu += point
	state := game.State

	state.LastSector = sector

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "score": point}).Info("Scored")

	state.Players[state.CurrentPlayer].Score -= point

	if state.Players[state.CurrentPlayer].Score > 0 {
		if game.doubleOut && state.Players[state.CurrentPlayer].Score == 1 {
			state.LastMsg = "You should end with a double"
			game.resetVisit()
			game.nextPlayer()
		} else {
			game.nextDart()
		}

	} else if state.Players[state.CurrentPlayer].Score == 0 {
		if game.doubleOut && sector.Pos != 2 {
			state.LastMsg = "You should end with a double"
			game.resetVisit()
			game.nextPlayer()
		} else {
			game.winner()
			if game.State.Ongoing == common.PLAYING {
				game.nextPlayer()
			}
		}

	} else {
		state.LastMsg = "You went beyond the target dude !"
		game.resetVisit()
		game.nextPlayer()
	}
	result = state
	return
}

func (game *Gamex01) winner() {
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

func (game *Gamex01) nextPlayer() {
	game.accu = 0
	state := game.State
	state.CurrentDart = 0
	state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Players)
	for state.Players[state.CurrentPlayer].Score == 0 {
		state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Players)
	}
	log.WithFields(log.Fields{"player": state.CurrentPlayer}).Info("Next player")
}

func (game *Gamex01) nextDart() {
	state := game.State
	if state.CurrentDart == 2 {
		game.nextPlayer()
	} else {
		state.CurrentDart += 1
		log.WithFields(log.Fields{"player": state.CurrentPlayer, "dart": state.CurrentDart}).Info("One more dart")
	}
}

func (game *Gamex01) resetVisit() {
	state := game.State
	state.Players[state.CurrentPlayer].Score += game.accu
}
