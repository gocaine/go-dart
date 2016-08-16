package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestGameCountupEnd(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd")

	game := NewGameCountUp(OptionCountUp{Target: 1})

	state := game.State()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer("test_board", "Bob")

	state = game.State()

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

	game := NewGameCountUp(OptionCountUp{Target: 301})

	state := game.State()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")

	AssertCurrents(t, game.state, 0, 0)

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 1)
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)

	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 1, 2)

	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)

	game.HoldOrNextPlayer()
	game.HoldOrNextPlayer()

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

func TestGameCountupOnHold(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupOnHold")

	game := NewGameCountUp(OptionCountUp{Target: 300})
	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")

	state, _ := game.HandleDart(common.Sector{Val: 5, Pos: 1})
	game.HoldOrNextPlayer()

	AssertGameState(t, state, common.ONHOLD)
	AssertCurrents(t, state, 0, 1)

	game.HoldOrNextPlayer()

	AssertGameState(t, state, common.PLAYING)
	AssertCurrents(t, state, 1, 0)

	game.HandleDart(common.Sector{Val: 5, Pos: 1})
	game.HandleDart(common.Sector{Val: 5, Pos: 1})
	game.HandleDart(common.Sector{Val: 5, Pos: 1})

	AssertGameState(t, state, common.ONHOLD)
	AssertCurrents(t, state, 1, 2)

	game.HoldOrNextPlayer()

	AssertGameState(t, state, common.PLAYING)
	AssertCurrents(t, state, 0, 0)
}
