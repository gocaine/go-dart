package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestGameHighestEnd(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameHighestEnd")

	game, err := NewGameHighest(map[string]interface{}{"Rounds": 0})

	expected := "Rounds should be at least 1"
	if err.Error() != expected {
		t.Errorf("Expected %s, but was %s", expected, err)
	}

	game, _ = NewGameHighest(map[string]interface{}{"Rounds": 1})

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

	game, _ := NewGameHighest(map[string]interface{}{"Rounds": 1})

	state := game.State()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 2})
	game.HandleDart(common.Sector{Val: 20, Pos: 1})

	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 2})
	game.HandleDart(common.Sector{Val: 20, Pos: 2})

	game.HoldOrNextPlayer()

	player := state.Players[0]

	AssertScore(t, player, 140)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	player = state.Players[1]

	AssertScore(t, player, 120)
	AssertRank(t, player, 2)
	AssertName(t, player, "Alice")
}

func TestGameHighestOnHold(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameHighestOnHold")

	game, _ := NewGameHighest(map[string]interface{}{"Rounds": 3})
	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")

	state, _ := game.HandleDart(common.Sector{Val: 5, Pos: 1})
	game.HoldOrNextPlayer()

	AssertGameState(t, state, common.ONHOLD)
	AssertCurrents(t, state, 0, 1)

	_, err := game.HandleDart(common.Sector{Val: 5, Pos: 1})
	AssertError(t, err, "Game is on hold and not ready to handle darts")

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
