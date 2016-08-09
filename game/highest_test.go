package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestGameHighestEnd(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameHighestEnd")

	game := NewGameHighest("test_board", OptionHighest{Rounds: 1})

	state := game.GetState()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer("Bob")

	state = game.GetState()

	if state.Ongoing != common.READY {
		t.Error("Game should be in ready mode")
	}

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})

	if state.Ongoing != common.OVER {
		t.Errorf("Game should be in OVER mode -- %+v", state)
	}

	score := state.Players[0]
	AssertScore(t, score, 180)
	AssertRank(t, score, 1)
}

func TestGameHighestEnd2Player(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameHighestEnd2Player")

	game := NewGameHighest("test_board", OptionHighest{Rounds: 1})

	state := game.GetState()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 2})
	game.HandleDart(common.Sector{Val: 20, Pos: 1})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 2})
	game.HandleDart(common.Sector{Val: 20, Pos: 2})

	player := state.Players[0]

	AssertScore(t, player, 140)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	player = state.Players[1]

	AssertScore(t, player, 120)
	AssertRank(t, player, 2)
	AssertName(t, player, "Alice")
}
