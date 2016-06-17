package common

type Sector struct {
	Val int
	Pos int
}

func (s Sector) IsValid() bool {
	if s.Val == 0 {
		return s.Pos == 0
	} else if s.Val > 0 && s.Val <= 20 {
		return s.Pos > 0 && s.Pos <= 3
	} else if s.Val == 25 {
		return s.Pos == 1 || s.Pos == 2
	}
	return false
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

type State int

const (
	INITIALIZING State = iota
	READY
	PLAYING
	OVER
)

type GameState struct {
	Scores        []Score
	Ongoing       State
	CurrentPlayer int
	CurrentDart   int
	LastMsg       string
	LastSector    Sector
}

func NewGameState() *GameState {

	g := new(GameState)
	g.Ongoing = INITIALIZING
	g.Scores = make([]Score, 0, 4)

	return g
}
