package common

type Sector struct {
	Val int
	Pos int
}

type Score struct {
	Player string
	Score  int
	Rank   int
}

// ByScore implements sort.Interface
type ByRank []Score

func (r ByRank) Len() int {
	return len(r)
}
func (r ByRank) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r ByRank) Less(i, j int) bool {
	if r[i].Rank == 0 {
		return r[j].Rank == 0
	} else if r[j].Rank == 0 {
		return true
	} else {
		return r[i].Rank < r[j].Rank
	}
}

type GameState struct {
	Scores        []Score
	Ongoing       bool
	CurrentPlayer int
	CurrentDart   int
}

func NewGameState() *GameState {

	g := new(GameState)

	g.Scores = make([]Score, 0, 4)

	return g
}
