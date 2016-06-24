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

type PlayerState struct {
	Name  string
	Score int
	Rank  int
	Histo map[string]int
}

// ByRank implements sort.Interface
type ByRank []PlayerState

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

// ByScore implements sort.Interface
type ByScore []PlayerState

func (r ByScore) Len() int {
	return len(r)
}
func (r ByScore) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r ByScore) Less(i, j int) bool {
	return r[i].Score > r[j].Score
}

type State int

const (
	INITIALIZING State = iota
	READY
	PLAYING
	OVER
)

type GameState struct {
	Players       []PlayerState
	Ongoing       State
	CurrentPlayer int
	CurrentDart   int
	LastMsg       string
	LastSector    Sector
	Round         int
}

func NewGameState() *GameState {

	g := new(GameState)
	g.Ongoing = INITIALIZING
	g.Players = make([]PlayerState, 0, 4)

	return g
}
