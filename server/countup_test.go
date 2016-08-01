package server

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestGameCountupEnd(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd")

	game := NewGameCountUp("test_board", OptionCountUp{Target: 1})

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

	if state.Ongoing != common.OVER {
		t.Errorf("Game should be in OVER mode -- %+v", state)
	}

	player := state.Players[0]

	if player.Score != 60 {
		t.Errorf("Player score should be 60 but was %d", player.Score)
	}

	if player.Rank != 1 {
		t.Errorf("Player rank should be 1 but was %d", player.Rank)
	}
}

func TestGameCountupEnd2Player(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd2Player")

	game := NewGameCountUp("test_board", OptionCountUp{Target: 301})

	state := game.GetState()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	AssertCurrents(t, game.State, 0, 0)

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.State, 0, 1)
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.State, 0, 2)
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.State, 1, 0)

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.State, 0, 0)

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, game.State, 1, 0)

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})

	player := state.Players[0]

	AssertScore(t, player, 360)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	player = state.Players[1]

	AssertScore(t, player, 300)
	AssertRank(t, player, 2)
	AssertName(t, player, "Alice")
}
