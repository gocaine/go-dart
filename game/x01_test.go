package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestGamex01End(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGamex01End")

	game := NewGamex01("test_board", Optionx01{Score: 1})
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	state, _ := game.HandleDart(common.Sector{Val: 5, Pos: 1})

	if state.Ongoing != common.PLAYING {
		t.Error("Game should not be ended")
	}

	alice := state.Players[0]

	if alice.Score != 1 {
		t.Error("Alice should have the same score : 1")
	}

	if state.CurrentPlayer != 1 || state.CurrentDart != 0 {
		t.Errorf("Should be bob's turn, first Dart (%d, %d)", state.CurrentPlayer, state.CurrentDart)
	}

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

	game := NewGamex01("test_board", Optionx01{Score: 501, DoubleOut: true})
	game.AddPlayer("Jack")

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 20, Pos: 3})

	game.HandleDart(common.Sector{Val: 20, Pos: 3})
	game.HandleDart(common.Sector{Val: 19, Pos: 3})
	state, _ := game.HandleDart(common.Sector{Val: 12, Pos: 2})

	if state.Ongoing != common.OVER {
		t.Error("Game should be Over")
	}

	ps := state.Players[0]

	AssertScore(t, ps, 0)
	AssertRank(t, ps, 1)
}

func TestGame301(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGame301")

	game := NewGamex01("test_board", Optionx01{Score: 301})
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.AddPlayer("Charly")
	game.AddPlayer("Dan")

	// Visit 1, Player 0
	state, _ := game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[0], 121)

	// Visit 1, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 2, 0)
	AssertScore(t, state.Players[1], 151)

	// Visit 1, Player 2
	state, _ = game.HandleDart(common.Sector{Val: 19, Pos: 2})
	AssertCurrents(t, state, 2, 1)
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 2, 2)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 3, 0)
	AssertScore(t, state.Players[2], 213)

	// Visit 1, Player 3
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 3, 1)
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 3, 2)
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 0, 0)
	AssertScore(t, state.Players[3], 301)

	// Visit 2, Player 0
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 7, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 2})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[0], 0)
	AssertRank(t, state.Players[0], 1)

	// Visit 2, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 1)
	AssertScore(t, state.Players[1], 91)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 2)
	AssertScore(t, state.Players[1], 31)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 2})
	AssertCurrents(t, state, 2, 0)
	AssertScore(t, state.Players[1], 151)

	// Visit 2, Player 2
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 2, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 2, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 0)
	AssertScore(t, state.Players[2], 33)

	// Visit 2, Player 3
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 3, 1)
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 3, 2)
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[3], 301)

	// Visit 3, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 1})
	AssertCurrents(t, state, 2, 0)
	AssertScore(t, state.Players[1], 11)

	// Visit 3, Player 2
	state, _ = game.HandleDart(common.Sector{Val: 10, Pos: 3})
	AssertCurrents(t, state, 2, 1)
	state, _ = game.HandleDart(common.Sector{Val: 1, Pos: 1})
	AssertCurrents(t, state, 2, 2)
	state, _ = game.HandleDart(common.Sector{Val: 1, Pos: 2})
	AssertCurrents(t, state, 3, 0)
	AssertScore(t, state.Players[2], 0)
	AssertRank(t, state.Players[2], 2)

	// Visit 3, Player 3
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[3], 121)

	// Visit 4, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 1, Pos: 3})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 5, Pos: 2})
	AssertCurrents(t, state, 3, 0)
	AssertScore(t, state.Players[1], 11)

	// Visit 4, Player 3
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 3, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[3], 121)

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

	game := NewGamex01("test_board", Optionx01{Score: 301, DoubleOut: true})
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	// Visit 1, Player 0
	state, _ := game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[0], 121)

	// Visit 1, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 25, Pos: 2})
	AssertCurrents(t, state, 0, 0)
	AssertScore(t, state.Players[1], 151)

	// Visit 2, Player 0
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[0], 121)

	// Visit 2, Player 1
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 1, 1)
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 1, 2)
	state, _ = game.HandleDart(common.Sector{Val: 0, Pos: 0})
	AssertCurrents(t, state, 0, 0)
	AssertScore(t, state.Players[1], 151)

	// Visit 3, Player 0
	state, _ = game.HandleDart(common.Sector{Val: 20, Pos: 3})
	AssertCurrents(t, state, 0, 1)
	state, _ = game.HandleDart(common.Sector{Val: 19, Pos: 3})
	AssertCurrents(t, state, 0, 2)
	state, _ = game.HandleDart(common.Sector{Val: 4, Pos: 1})
	AssertCurrents(t, state, 1, 0)
	AssertScore(t, state.Players[0], 121)

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
