package server

import "go-dart/common"

type Gamex01 struct {
	score   int
	started bool
	players []string
}

func NewGame(score int) *Gamex01 {
	g := new(Gamex01)

	g.score = score
	g.players = make([]string, 4)

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
	if !game.started {
		game.started = true
	} else {
		panic("Game already started")
	}
}

func (game *Gamex01) HandleDart(sector common.Sector) common.GameState {
	
}

