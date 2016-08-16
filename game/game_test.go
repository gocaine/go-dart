package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestAddPlayer(t *testing.T) {
	fmt.Println()
	fmt.Println("TestAddPlayer")

	game := new(AGame)
	game.state = common.NewGameState()

	err := game.Start()
	if err == nil {
		fatalStack(t, "should not be able to start a game w/o players")
	}

	err = game.AddPlayer("b1", "p1")
	if err != nil {
		fatalStack(t, "should be able to add a player")
	}

	err = game.AddPlayer("b1", "p1")
	if err == nil {
		fatalStack(t, "should not be able to add twice a player")
	}

	err = game.AddPlayer("b2", "p2")
	if err != nil {
		fatalStack(t, "should be able to add a player")
	}

	err = game.Start()
	if err != nil {
		fatalStack(t, "should be able to start a game w/ players")
	}
}

func TestBoardHasLeft(t *testing.T) {
	fmt.Println()
	fmt.Println("TestBoardHasLeft")

	game := new(AGame)
	game.state = common.NewGameState()

	game.Start()
	game.AddPlayer("b1", "p1")
	game.AddPlayer("b2", "p2")
	game.Start()
	if !game.BoardHasLeft("another board") {
		fatalStack(t, "should not be concerned")
	}

	if game.BoardHasLeft("b2") {
		fatalStack(t, "should have ended the game")
	}

	if game.state.Ongoing != common.OVER {
		fatalStack(t, "game should be over")
	}
}
