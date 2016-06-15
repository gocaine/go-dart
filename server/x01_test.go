package server

import (
	"fmt"
	"go-dart/common"
	"testing"
)

func TestGamex01End(t *testing.T) {
	game := NewGamex01(1)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.Start()
	state := game.HandleDart(common.Sector{Val: 5, Pos: 1})

	if !state.Ongoing {
		t.Error("Game should not be ended")
	}

	alice := state.Scores[0]

	if alice.Score != 1 {
		t.Error("Alice should have the same score : 1")
	}

	if state.CurrentPlayer != 1 || state.CurrentDart != 0 {
		t.Errorf("Should be bob's turn, first Dart (%d, %d)", state.CurrentPlayer, state.CurrentDart)
	}

	state = game.HandleDart(common.Sector{Val: 1, Pos: 1})

	if state.Ongoing {
		t.Error("Game should be ended")
	}

	bob := state.Scores[0]

	if bob.Player != "Bob" {
		t.Error("Bob should be first")
	}

	fmt.Printf("State is %+v \n", state)

}
