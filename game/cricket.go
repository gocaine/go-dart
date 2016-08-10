package game

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
)

var SECTORS = [...]string{"15", "16", "17", "18", "19", "20", "25"}

type GameCricket struct {
	AGame
	noScore   bool
	cutThroat bool
	memory    map[string]int
}

type OptionCricket struct {
	NoScore   bool
	CutThroat bool
}

func NewGameCricket(board string, opt OptionCricket) *GameCricket {

	g := new(GameCricket)
	g.SetBoard(board)
	g.noScore = opt.NoScore
	g.cutThroat = opt.CutThroat
	g.State = common.NewGameState()
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

func (game *GameCricket) AddPlayer(name string) (error error) {

	error = game.AGame.AddPlayer(name)
	if error == nil {
		game.State.Players[len(game.State.Players)-1].Histo = make(map[string]int)
	}
	log.WithFields(log.Fields{"name": name, "player": game.State.Players[len(game.State.Players)-1]}).Info("AddPlayer")
	return
}

func (game *GameCricket) Start() (error error) {

	error = game.AGame.Start()
	if error == nil {
		for _, key := range SECTORS {
			game.memory[key] = len(game.State.Players)
		}
		log.WithFields(log.Fields{"memory": game.memory}).Info("Start")
	}
	return
}

func (game *GameCricket) HandleDart(sector common.Sector) (result *common.GameState, error error) {

	if game.State.Ongoing == common.READY {
		// first dart starts the game
		err := game.Start()
		if err != nil {
			error = err
			return
		}
	}

	if game.State.Ongoing != common.PLAYING {
		error = errors.New("Game is not started or is ended")
		return
	}

	if !sector.IsValid() {
		log.WithFields(log.Fields{"sector": sector}).Error("Invalid sector")
		error = errors.New("Sector is not a valid one")
		return
	}

	state := game.State

	state.LastSector = sector
	sVal := strconv.Itoa(sector.Val)

	log.WithFields(log.Fields{"player": state.CurrentPlayer, "sector": sector}).Info("Hit")

	if sector.Val >= 15 {
		var count int = state.Players[state.CurrentPlayer].Histo[sVal]
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

func (game *GameCricket) score(val, pos int) {
	log.WithFields(log.Fields{"sector": val, "number": pos}).Info("score")
	if game.noScore {
		// no score at all
	} else {
		points := val * pos
		game.State.LastMsg = fmt.Sprintf("Scored : %d", points)
		if game.cutThroat {
			for key := range game.State.Players {
				if game.State.Players[key].Histo[strconv.Itoa(val)] < 3 {
					game.State.Players[key].Score += points
				}
			}
		} else {
			game.State.Players[game.State.CurrentPlayer].Score += points
		}
	}

	game.checkWinner()
}

func (game *GameCricket) checkWinner() {
	log.WithFields(log.Fields{"state": game.State}).Info("checkWinner")
	player := game.State.Players[game.State.CurrentPlayer]
	remain := false
	for key := 0; key < len(SECTORS) && !remain; key++ {
		remain = player.Histo[SECTORS[key]] != 3
	}
	// The player has opened everything if for none of the sector hits are missing
	if !remain {
		// if we are in noScore mode, no more hit remaining is a sufficient condition
		if game.noScore {
			game.winner()
		} else {
			if game.cutThroat {
				if lowest(game.State.Players, game.State.CurrentPlayer) {
					game.winner()
				} else {
					game.nextDart()
				}
			} else {
				if highest(game.State.Players, game.State.CurrentPlayer) {
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

func (game *GameCricket) winner() {
	log.Info("winner")
	state := game.State
	game.State.LastMsg = fmt.Sprintf("Winner : %s", state.Players[state.CurrentPlayer].Name)
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	state.LastMsg = fmt.Sprintf("Player %d end at rank #%d", state.CurrentPlayer, game.rank+1)
	game.rank++
	if game.rank >= len(state.Players)-1 {
		game.State.Ongoing = common.OVER
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
