package game

import (
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/i18n"
	"log"
	"runtime"
	"testing"
)

// AssertScore assert the score is what expected
func AssertScore(t *testing.T, ps common.PlayerState, target int) {
	if ps.Score != target {
		fatalStack(t, "Player score should be %d but was %d -- %+v", target, ps.Score, ps)
	}
}

// AssertRank assert the ranking is what expected
func AssertRank(t *testing.T, ps common.PlayerState, target int) {
	if ps.Rank != target {
		fatalStack(t, "Player rank should be %d but was %d", target, ps.Rank)
	}
}

// AssertName assert the name is what expected
func AssertName(t *testing.T, ps common.PlayerState, name string) {
	if ps.Name != name {
		fatalStack(t, "Player name should be %s but was %s", name, ps.Name)
	}
}

// AssertCurrents assertion on currentPlayer and CurrentDart
func AssertCurrents(t *testing.T, state *common.GameState, p, d int) {
	if state.CurrentPlayer != p || state.CurrentDart != d {
		fatalStack(t, "Player should be %d and Dart %d, but was %d and %d -- %+v", p, d, state.CurrentPlayer, state.CurrentDart, state)
	}
}

// AssertEquals classical assertEquals
func AssertEquals(t *testing.T, expected, actual interface{}) {
	if actual != expected {
		fatalStack(t, "Expected : %+v, Was : %+v", expected, actual)
	}
}

// AssertGameState assert the GameState is as exepcetd
func AssertGameState(t *testing.T, state *common.GameState, expected common.State) {
	if state.Ongoing != expected {
		fatalStack(t, "Game state should be %d but was %d -- %+v", expected, state.Ongoing, state)
	}

}

// AssertError assert the error content
func AssertError(t *testing.T, err error, expected string) {

	if value, ok := i18n.BaseTranslation(err.Error()); !ok || value != expected {
		fatalStack(t, "Error should be %s but was %s", expected, value)
	}

}

func fatalStack(t *testing.T, format string, args ...interface{}) {
	//stack trace
	var stack [4096]byte
	runtime.Stack(stack[:], false)
	log.Printf("%s\n", stack[:])

	//fatal error
	t.Fatalf(format, args...)
}
