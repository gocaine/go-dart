package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestGamex01End(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGamex01End")

	game, err := NewGamex01(map[string]interface{}{"Score": "aa"})

	expected := "aa is an invalid value for Score"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected %s, but was %s", expected, err)
	}

	game, err = NewGamex01(map[string]interface{}{"Score": 1})

	expected = "Score should be at least 61"
	if err.Error() != expected {
		t.Errorf("Expected %s, but was %s", expected, err)
	}

	game, _ = NewGamex01(map[string]interface{}{"Score": 61})
	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")
	state, _ := game.HandleDart(common.Sector{Val: 20, Pos: 3})
	state, _ = game.HandleDart(common.Sector{Val: 5, Pos: 1})

	AssertGameState(t, state, common.ONHOLD)

	game.HoldOrNextPlayer()
	AssertGameState(t, state, common.PLAYING)

	alice := state.Players[0]

	AssertScore(t, alice, 61)

	if state.CurrentPlayer != 1 || state.CurrentDart != 0 {
		t.Errorf("Should be bob's turn, first Dart (%d, %d)", state.CurrentPlayer, state.CurrentDart)
	}
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	state, _ = game.HandleDart(common.Sector{Val: 1, Pos: 1})

	if state.Ongoing != common.OVER {
		t.Error("Game should be ended")
	}

	bob := state.Players[0]

	if bob.Name != "Bob" {
		t.Error("Bob should be first")
	}

}

func TestGamex01SoloEnd(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGamex01SoloEnd")

	game, _ := NewGamex01(map[string]interface{}{"Score": 501, "DoubleOut": true})
	game.AddPlayer("test_board", "Jack")

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})

	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})

	game.HoldOrNextPlayer()

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 19, Pos: 3})
	state, _ := game.HandleDart(common.Sector{Val: 12, Pos: 2})

	AssertGameState(t, state, common.OVER)

	ps := state.Players[0]

	AssertScore(t, ps, 0)
	AssertRank(t, ps, 1)
}

func TestGame301(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGame301")

	game, _ := NewGamex01(map[string]interface{}{"Score": 301})
	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")
	game.AddPlayer("test_board", "Charly")
	game.AddPlayer("test_board", "Dan")

	// Visit 1, Player 0
	state, _ := game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	AssertScore(t, state.Players[0], 121)

	game.HoldOrNextPlayer()

	// Visit 1, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 2)
	AssertScore(t, state.Players[1], 151)

	game.HoldOrNextPlayer()

	// Visit 1, Player 2
	state, _ = game.HandleDart(common.Sector{Val: 19, Pos: 2})
	AssertCurrents(t, state, 2, 1)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 2, 2)
	game.HoldOrNextPlayer()
	AssertScore(t, state.Players[2], 213)

	game.HoldOrNextPlayer()

	// Visit 1, Player 3
	game.HoldOrNextPlayer()
	AssertCurrents(t, state, 3, 0)
	AssertScore(t, state.Players[3], 301)

	game.HoldOrNextPlayer()

	// Visit 2, Player 0
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 7, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 2})
	AssertCurrents(t, state, 0, 2)
	AssertScore(t, state.Players[0], 0)
	AssertRank(t, state.Players[0], 1)

	game.HoldOrNextPlayer()

	// Visit 2, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 1)
	AssertScore(t, state.Players[1], 91)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 2)
	AssertScore(t, state.Players[1], 31)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 2})
	AssertCurrents(t, state, 1, 2)
	AssertScore(t, state.Players[1], 151)

	game.HoldOrNextPlayer()

	// Visit 2, Player 2
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 2, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 2, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 2, 2)
	AssertScore(t, state.Players[2], 33)

	game.HoldOrNextPlayer()

	// Visit 2, Player 3
	game.HoldOrNextPlayer()
	AssertCurrents(t, state, 3, 0)
	AssertScore(t, state.Players[3], 301)

	game.HoldOrNextPlayer()

	// Visit 3, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 1})
	AssertCurrents(t, state, 1, 2)
	AssertScore(t, state.Players[1], 11)

	game.HoldOrNextPlayer()

	// Visit 3, Player 2
	state, _ = game.HandleDart(common.Sector{Val: 10, Pos: 3})
	AssertCurrents(t, state, 2, 1)
	state, _ = game.HandleDart(common.Sector{Val: 1, Pos: 1})
	AssertCurrents(t, state, 2, 2)
	state, _ = game.HandleDart(common.Sector{Val: 1, Pos: 2})
	AssertCurrents(t, state, 2, 2)
	AssertScore(t, state.Players[2], 0)
	AssertRank(t, state.Players[2], 2)

	game.HoldOrNextPlayer()

	// Visit 3, Player 3
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 2)
	AssertScore(t, state.Players[3], 121)

	game.HoldOrNextPlayer()

	// Visit 4, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 1, Pos: 3})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 5, Pos: 2})
	AssertCurrents(t, state, 1, 1)
	AssertScore(t, state.Players[1], 11)

	game.HoldOrNextPlayer()

	// Visit 4, Player 3
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 2)
	AssertScore(t, state.Players[3], 121)

	game.HoldOrNextPlayer()

	// Visit 5, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 3, Pos: 2})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 5, Pos: 1})

	if state.Ongoing != common.OVER {
		t.Error("Game should be ended")
	}

	AssertScore(t, state.Players[0], 0)
	AssertScore(t, state.Players[1], 0)
	AssertScore(t, state.Players[2], 0)
	AssertScore(t, state.Players[3], 121)

	AssertRank(t, state.Players[0], 1)
	AssertRank(t, state.Players[1], 2)
	AssertRank(t, state.Players[2], 3)
	AssertRank(t, state.Players[3], 4)

	AssertName(t, state.Players[0], "Alice")
	AssertName(t, state.Players[1], "Charly")
	AssertName(t, state.Players[2], "Bob")
	AssertName(t, state.Players[3], "Dan")

}

func TestGame301DoubleOut(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGame301DoubleOut")

	game, _ := NewGamex01(map[string]interface{}{"Score": 301, "DoubleOut": true})
	game.AddPlayer("test_board", "Alice")
	game.AddPlayer("test_board", "Bob")

	// Visit 1, Player 0
	state, _ := game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	AssertScore(t, state.Players[0], 121)

	game.HoldOrNextPlayer()

	// Visit 1, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 2)
	AssertScore(t, state.Players[1], 151)

	game.HoldOrNextPlayer()

	// Visit 2, Player 0
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	AssertScore(t, state.Players[0], 121)

	game.HoldOrNextPlayer()

	// Visit 2, Player 1
	game.HoldOrNextPlayer()
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[1], 151)

	game.HoldOrNextPlayer()

	// Visit 3, Player 0
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 19, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 4, Pos: 1})
	AssertCurrents(t, state, 0, 2)
	AssertScore(t, state.Players[0], 121)

	game.HoldOrNextPlayer()

	// Visit 3, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 19, Pos: 3})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 17, Pos: 2})

	if state.Ongoing != common.OVER {
		t.Error("Game should be ended")
	}

	AssertScore(t, state.Players[0], 0)
	AssertScore(t, state.Players[1], 121)

	AssertRank(t, state.Players[0], 1)
	AssertRank(t, state.Players[1], 2)

	AssertName(t, state.Players[0], "Bob")
	AssertName(t, state.Players[1], "Alice")

}

func TestGamex01OnHold(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGamex01OnHold")

	game, _ := NewGamex01(map[string]interface{}{"Score": 301})
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
