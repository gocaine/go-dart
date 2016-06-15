package common

type Sector struct {
	Val int
	Pos int
}

type Score struct {
	Player string
	Score  int
}

type GameState struct {
	Scores        []Score
	Ongoing       bool
	CurrentPlayer int
	CurrentDart   int
}

func NewGameState() GameState {
	g := GameState{Scores: make([]Score, 0, 4)}

	return g
}
