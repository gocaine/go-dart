package game

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
)

var sectors = [...]string{"15", "16", "17", "18", "19", "20", "25"}

// Cricket is a cricket series Game (Cricket, Cut-throat)
type Cricket struct {
	AGame
	noScore   bool
	cutThroat bool
	memory    map[string]int
}

// OptionCricket is the struct to handle GameCricket parameters
type OptionCricket struct {
	NoScore   bool
	CutThroat bool
}

// NewGameCricket : GameCricket constructor using a OptionCricket
func NewGameCricket(board string, opt OptionCricket) *Cricket {

	g := new(Cricket)
	g.SetBoard(board)
	g.noScore = opt.NoScore
	g.cutThroat = opt.CutThroat
	g.state = common.NewGameState()
	dStyle := "Cricket"
	if opt.CutThroat {
		dStyle = "Cut-Throat Cricket"
	} else if opt.NoScore {
		dStyle = "No Score Cricket"
	}
	g.DisplayStyle = dStyle
	g.memory = make(map[string]int)

	return g
}

// AddPlayer add a new player to the game
func (game *Cricket) AddPlayer(name string) (error error) {

	error = game.AGame.AddPlayer(name)
	if error == nil {
		game.state.Players[len(game.state.Players)-1].Histo = make(map[string]int)
	}
	log.WithFields(log.Fields{"name": name, "player": game.state.Players[len(game.state.Players)-1]}).Info("AddPlayer")
	return
}

// Start start the game, Darts will be handled
func (game *Cricket) Start() (error error) {

	error = game.AGame.Start()
	if error == nil {
		for _, key := range sectors {
			game.memory[key] = len(game.state.Players)
		}
		log.WithFields(log.Fields{"memory": game.memory}).Info("Start")
	}
	return
}

// HandleDart the implementation has to handle the Dart regarding the current player, the cricket rules, and the context. Return a GameState
func (game *Cricket) HandleDart(sector common.Sector) (result *common.GameState, error error) {

	if game.state.Ongoing == common.READY {
		// first dart starts the game
		err := game.Start()
		if err != nil {
			error = err
			return
		}
	}

	if game.state.Ongoing != common.PLAYING {
		error = errors.New("Game is not started or is ended")
		return
	}

	if !sector.IsValid() {
		log.WithFields(log.Fields{"sector": sector}).Error("Invalid sector")
		error = errors.New("Sector is not a valid one")
		return
	}

	state := game.state

	state.LastSector = sector
	sVal := strconv.Itoa(sector.Val)

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "sector": sector}).Info("Hit")

	if sector.Val >= 15 {
		var count = state.Players[state.CurrentPlayer].Histo[sVal]
		log.WithFields(log.Fields{"count": count, "val": sector.Val}).Info("Hit Count")
		if count == 3 {
			open := game.memory[sVal] > 0
			if open {
				game.score(sector.Val, sector.Pos)
			} else {
				game.nextDart()
			}
		} else {
			remain := 0
			count = count + sector.Pos
			if count > 3 {
				remain = count - 3
				count = 3
			}
			state.Players[state.CurrentPlayer].Histo[sVal] = count
			if count == 3 {
				game.memory[sVal] = game.memory[sVal] - 1
				open := game.memory[sVal] > 0
				if open {
					state.LastMsg = fmt.Sprintf("Opened : %s", sVal)
				} else {
					state.LastMsg = fmt.Sprintf("Closed : %s", sVal)
				}
				if open && remain > 0 {
					game.score(sector.Val, remain)
				} else {
					game.checkWinner()
				}
			} else {
				state.LastMsg = fmt.Sprintf("Hit : %d x %s", sector.Pos, sVal)
				game.nextDart()
			}
		}
	} else {
		game.nextDart()
	}
	result = state
	return
}

func (game *Cricket) score(val, pos int) {
	log.WithFields(log.Fields{"sector": val, "number": pos}).Info("score")
	if game.noScore {
		// no score at all
	} else {
		points := val * pos
		game.state.LastMsg = fmt.Sprintf("Scored : %d", points)
		if game.cutThroat {
			for key := range game.state.Players {
				if game.state.Players[key].Histo[strconv.Itoa(val)] < 3 {
					game.state.Players[key].Score += points
				}
			}
		} else {
			game.state.Players[game.state.CurrentPlayer].Score += points
		}
	}

	game.checkWinner()
}

func (game *Cricket) checkWinner() {
	log.WithFields(log.Fields{"state": game.state}).Info("checkWinner")
	player := game.state.Players[game.state.CurrentPlayer]
	remain := false
	for key := 0; key < len(sectors) && !remain; key++ {
		remain = player.Histo[sectors[key]] != 3
	}
	// The player has opened everything if for none of the sector hits are missing
	if !remain {
		// if we are in noScore mode, no more hit remaining is a sufficient condition
		if game.noScore {
			game.winner()
		} else {
			if game.cutThroat {
				if lowest(game.state.Players, game.state.CurrentPlayer) {
					game.winner()
				} else {
					game.nextDart()
				}
			} else {
				if highest(game.state.Players, game.state.CurrentPlayer) {
					game.winner()
				} else {
					game.nextDart()
				}
			}
		}
	} else {
		game.nextDart()
	}
}

func (game *Cricket) winner() {
	log.Info("winner")
	state := game.state
	game.state.LastMsg = fmt.Sprintf("Winner : %s", state.Players[state.CurrentPlayer].Name)
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	state.LastMsg = fmt.Sprintf("Player %d end at rank #%d", state.CurrentPlayer, game.rank+1)
	game.rank++
	if game.rank >= len(state.Players)-1 {
		game.state.Ongoing = common.OVER
		sort.Sort(common.ByRank(state.Players))
		if len(state.Players) > 1 {
			state.Players[len(state.Players)-1].Rank = game.rank + 1
		}
	} else {
		game.nextPlayer()
	}
}

// check if current player has the highest score
func highest(players []common.PlayerState, current int) bool {
	target := players[current].Score
	for key, val := range players {
		if key != current && val.Rank == 0 && val.Score > target {
			return false
		}
	}
	return true
}

// check if current player has the lowest score
func lowest(players []common.PlayerState, current int) bool {
	target := players[current].Score
	for key, val := range players {
		if key != current && val.Rank == 0 && val.Score < target {
			return false
		}
	}
	return true
}
