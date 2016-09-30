package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/i18n"
)

func TestGameCountupI18N(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupI18N")

	_, err := NewGameCountUp(map[string]interface{}{"Target": 1})

	expected := "Target should be at least 61"
	value, ok := i18n.BaseTranslation(err.Error())
	if !ok || value != expected {
		t.Errorf("Expected %s, but was %s", expected, value)
	}
}

func TestGameCountupEnd(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd")

	game, _ := NewGameCountUp(map[string]interface{}{"Target": 61})

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

	game, _ := NewGameCountUp(map[string]interface{}{"Target": 301})

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

func TestGameCountupEnd3Player(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCountupEnd3Player")

	game, _ := NewGameCountUp(map[string]interface{}{"Target": 301})

	state := game.State()

	if state.Ongoing != common.INITIALIZING {
		t.Error("Game should be in initializing mode")
	}

	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")
	game.AddPlayer("test_board", "Charly")

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
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 2, 2)

	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 0, 2)

	game.HoldOrNextPlayer()
	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, game.state, 1, 2)

	AssertGameState(t, state, common.ONHOLD)

	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})

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

	game, _ := NewGameCountUp(map[string]interface{}{"Target": 300})
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
