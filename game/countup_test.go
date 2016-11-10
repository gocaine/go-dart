package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestGameCountupEnd(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd")

	ctx := createContext("eng")

	game, err := NewGameCountUp(ctx, map[string]interface{}{"Target": 1})

	expected := "Target should be at least 61"
	if err.Error() != expected {
		t.Errorf("Expected %s, but was %s", expected, err)
	}

	game, _ = NewGameCountUp(ctx, map[string]interface{}{"Target": 61})

	state := game.State()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer(ctx, "test_board", "Bob")

	state = game.State()

	if state.Ongoing != common.READY {
		t.Error("Game should be in ready mode")
	}

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})

	if state.Ongoing != common.OVER {
		t.Errorf("Game should be in OVER mode -- %+v", state)
	}

	player := state.Players[0]

	AssertScore(t, player, 120)
	AssertRank(t, player, 1)
}

func TestGameCountupEnd2Player(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd2Player")

	ctx := createContext("eng")

	game, _ := NewGameCountUp(ctx, map[string]interface{}{"Target": 301})

	state := game.State()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer(ctx, "test_board", "Alice")
	err := game.AddPlayer(ctx, "test_board", "Alice")
	AssertError(t, err, "Player name is already in use")

	game.AddPlayer(ctx, "test_board", "Bob")

	AssertCurrents(t, game.state, 0, 0)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 1)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)

	game.HoldOrNextPlayer(ctx)

	err = game.AddPlayer(ctx, "test_board", "Jack")
	AssertError(t, err, "Player cannot be added")

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 1, 2)

	game.HoldOrNextPlayer(ctx)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)

	game.HoldOrNextPlayer(ctx)
	game.HoldOrNextPlayer(ctx)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})

	player := state.Players[0]

	AssertScore(t, player, 360)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	player = state.Players[1]

	AssertScore(t, player, 300)
	AssertRank(t, player, 2)
	AssertName(t, player, "Alice")
}

func TestGameCountupEnd3Player(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd3Player")

	ctx := createContext("eng")

	game, _ := NewGameCountUp(ctx, map[string]interface{}{"Target": 301})

	state := game.State()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer(ctx, "test_board", "Alice")
	game.AddPlayer(ctx, "test_board", "Bob")
	game.AddPlayer(ctx, "test_board", "Charly")

	AssertCurrents(t, game.state, 0, 0)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 1)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)

	game.HoldOrNextPlayer(ctx)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 1, 2)

	game.HoldOrNextPlayer(ctx)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 2, 2)

	game.HoldOrNextPlayer(ctx)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)

	game.HoldOrNextPlayer(ctx)
	game.HoldOrNextPlayer(ctx)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 1, 2)

	AssertGameState(t, state, common.ONHOLD)

	game.HoldOrNextPlayer(ctx)

	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})

	player := state.Players[0]

	AssertScore(t, player, 360)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	player = state.Players[1]

	AssertScore(t, player, 360)
	AssertRank(t, player, 2)
	AssertName(t, player, "Charly")

	player = state.Players[2]

	AssertScore(t, player, 300)
	AssertRank(t, player, 3)
	AssertName(t, player, "Alice")
}

func TestGameCountupOnHold(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupOnHold")

	ctx := createContext("eng")

	game, _ := NewGameCountUp(ctx, map[string]interface{}{"Target": 300})
	game.AddPlayer(ctx, "test_board", "Alice")
	game.AddPlayer(ctx, "test_board", "Bob")

	state, _ := game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.ONHOLD)
	AssertCurrents(t, state, 0, 1)

	_, err := game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	AssertError(t, err, "Game is on hold and not ready to handle darts")

	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.PLAYING)
	AssertCurrents(t, state, 1, 0)

	game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})

	AssertGameState(t, state, common.ONHOLD)
	AssertCurrents(t, state, 1, 2)

	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.PLAYING)
	AssertCurrents(t, state, 0, 0)
}
