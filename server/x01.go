package server

import "go-dart/common"

type Gamex01 struct {
	score   int
	started bool
	players []string
	remains []int
}

func NewGamex01(score int) *Gamex01 {
	g := new(Gamex01)

	g.score = score
	g.players = make([]string, 0, 4)

	return g
}

func (game *Gamex01) AddPlayer(name string) {
	if !game.started {
		game.players = append(game.players, name)
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) Start() {
	if !game.started && len(game.players) > 0 {
		game.started = true
		game.remains = make([]int, len(game.players))
		for i := range game.remains {
			game.remains[i] = game.score
		}
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) HandleDart(sector common.Sector) common.GameState {

	return {}
}
