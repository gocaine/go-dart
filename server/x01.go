package server

import (
	"go-dart/common"
	"sort"
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
	if !game.State.Ongoing {
		game.State.Scores = append(game.State.Scores, common.Score{Player: name, Score: game.score})
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) Start() {
	if !game.State.Ongoing && len(game.State.Scores) > 0 && game.score > 0 {
		state := game.State
		state.Ongoing = true
		state.CurrentPlayer = 0
		state.CurrentDart = 0
		for i := range state.Scores {
			state.Scores[i].Score = game.score
		}
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) HandleDart(sector common.Sector) *common.GameState {

	if !game.State.Ongoing {
		panic("Game is not started or is ended")
	}

	if !sector.IsValid() {
		panic("Sector is not a valid one")
	}

	point := int(sector.Val) * int(sector.Pos)
	game.accu += point
	state := game.State
	state.Scores[state.CurrentPlayer].Score -= point

	if state.Scores[state.CurrentPlayer].Score > 0 {
		game.nextDart()

	} else if state.Scores[state.CurrentPlayer].Score == 0 {
		game.winner()
		if game.State.Ongoing {
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
	if game.rank == len(state.Scores)-1 {
		state.Ongoing = false
		sort.Sort(common.ByRank(state.Scores))
		state.Scores[len(state.Scores)-1].Rank = game.rank + 1
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
}

func (game *Gamex01) nextDart() {
	state := game.State
	if state.CurrentDart == 2 {
		game.nextPlayer()
	} else {
		state.CurrentDart += 1
	}
}

func (game *Gamex01) resetVisit() {
	state := game.State
	state.Scores[state.CurrentPlayer].Score += game.accu
}
