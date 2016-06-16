package server

import (
	"go-dart/common"
	"sort"
	log "github.com/Sirupsen/logrus"
)

type Gamex01 struct {
	score int
	State *common.GameState
	accu  int
	rank  int
}

func NewGamex01(score int) *Gamex01 {
	g := new(Gamex01)

	g.score = score
	g.State = common.NewGameState()

	return g
}

func (game *Gamex01) AddPlayer(name string) {
	if game.State.Ongoing == common.INITIALIZING || game.State.Ongoing == common.READY {
		log.WithFields(log.Fields{"player": name}).Infof("Player added to the game", name)
		game.State.Scores = append(game.State.Scores, common.Score{Player: name, Score: game.score})
		// now that we have at least one player, we are in a ready state, waiting for other players or the first dart
		game.State.Ongoing = common.READY
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) Start() {
	if game.State.Ongoing == common.READY && len(game.State.Scores) > 0 && game.score > 0 {
		state := game.State
		state.Ongoing = common.PLAYING
		state.CurrentPlayer = 0
		state.CurrentDart = 0
		for i := range state.Scores {
			state.Scores[i].Score = game.score
		}
		log.Infof("The game is now started")
	} else {
		panic("Game cannot start")
	}
}

func (game *Gamex01) HandleDart(sector common.Sector) *common.GameState {

	if game.State.Ongoing == common.READY {
		// first dart starts the game
		game.Start()
	}

	if game.State.Ongoing != common.PLAYING {
		panic("Game is not started or is ended")
	}

	if !sector.IsValid() {
		log.WithFields(log.Fields{"sector": sector}).Error("Invalid sector")
		panic("Sector is not a valid one")
	}

	point := int(sector.Val) * int(sector.Pos)
	game.accu += point
	state := game.State

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "score": point}).Info("Scored")

	state.Scores[state.CurrentPlayer].Score -= point

	if state.Scores[state.CurrentPlayer].Score > 0 {
		game.nextDart()

	} else if state.Scores[state.CurrentPlayer].Score == 0 {
		game.winner()
		if game.State.Ongoing == common.PLAYING {
			game.nextPlayer()
		}

	} else {
		game.resetVisit()
		game.nextPlayer()
	}

	return state
}

func (game *Gamex01) winner() {
	state := game.State
	state.Scores[state.CurrentPlayer].Rank = game.rank + 1
	game.rank++
	if game.rank == len(state.Scores) - 1 {
		game.State.Ongoing = common.OVER
		sort.Sort(common.ByRank(state.Scores))
		state.Scores[len(state.Scores) - 1].Rank = game.rank + 1
	}
}

func (game *Gamex01) nextPlayer() {
	game.accu = 0
	state := game.State
	state.CurrentDart = 0
	state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Scores)
	for state.Scores[state.CurrentPlayer].Score == 0 {
		state.CurrentPlayer = (state.CurrentPlayer + 1) % len(state.Scores)
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
	state.Scores[state.CurrentPlayer].Score += game.accu
}
