package server

import "go-dart/common"

type Gamex01 struct {
	score int
	State common.GameState
	accu  int
}

func NewGamex01(score int) *Gamex01 {
	g := new(Gamex01)

	g.score = score
	g.State = common.NewGameState()

	return g
}

func (game *Gamex01) AddPlayer(name string) {
	if !game.started {
		append(game.players, name)
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) Start() {
	if !game.State.Ongoing && len(game.State.Scores) > 0 {
		state := game.State
		state.Ongoing = true
		state.CurrentPlayer = 0
		state.CurrentDart = 0
		for i := range state.Scores {
			state.Scores[i] = game.score
		}
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) HandleDart(sector common.Sector) common.GameState {

	point := sector.Name * sector.Pos
	game.accu += point
	state := game.State
	if game.accu < state.Scores[state.CurrentPlayer].Score {
		
	} else {
		game.nextPlayer()
	}
	return common.GameState{}
}

