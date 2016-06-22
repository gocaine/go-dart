package server

import (
	"go-dart/common"
	"log"
	"runtime"
	"testing"
)

func AssertScore(t *testing.T, ps common.PlayerState, target int) {
	if ps.Score != target {
		fatalStack(t, "Player score should be %d but was %d", target, ps.Score)
	}
}

func AssertRank(t *testing.T, ps common.PlayerState, target int) {
	if ps.Rank != target {
		fatalStack(t, "Player rank should be %d but was %d", target, ps.Rank)
	}
}

func AssertName(t *testing.T, ps common.PlayerState, name string) {
	if ps.Name != name {
		fatalStack(t, "Player name should be %s but was %s", name, ps.Name)
	}
}

func AssertCurrents(t *testing.T, state *common.GameState, p, d int) {
	if state.CurrentPlayer != p || state.CurrentDart != d {
		fatalStack(t, "Player should be %d and Dart %d, but was %d and %d -- %+v", p, d, state.CurrentPlayer, state.CurrentDart, state)
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
