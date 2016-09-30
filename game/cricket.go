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
	BaseGame
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
func NewGameCricket(opts map[string]interface{}) (g *Cricket, err error) {
	opt := newOptionCricket(opts)
	if opt.CutThroat && opt.NoScore {
		err = errors.New("game.cricket.error.incompatible")
		return
	}
	g = new(Cricket)
	g.noScore = opt.NoScore
	g.cutThroat = opt.CutThroat
	g.state = common.NewGameState()
	dStyle := "game.cricket.display.cricket"
	if opt.CutThroat {
		dStyle = "game.cricket.display.cutthroat"
	} else if opt.NoScore {
		dStyle = "game.cricket.display.noscore"
	}
	g.DisplayStyle = dStyle
	g.memory = make(map[string]int)

	return
}

// AddPlayer add a new player to the game
func (game *Cricket) AddPlayer(board string, name string) (error error) {

	error = commonAddPlayer(game, board, name)
	if error == nil {
		game.state.Players[len(game.state.Players)-1].Histo = make(map[string]int)
	}
	log.WithFields(log.Fields{"name": name, "player": game.state.Players[len(game.state.Players)-1]}).Info("AddPlayer")
	return
}

// Start start the game, Darts will be handled
func (game *Cricket) Start() (error error) {

	error = commonStart(game)
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

	error = commonHandleDartChecks(game, sector)
	if error != nil {
		return
	}

	state := game.state

	state.LastSector = sector
	sVal := strconv.Itoa(sector.Val)

	state.Players[state.CurrentPlayer].Visits = append(state.Players[state.CurrentPlayer].Visits, sector)

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
					state.LastMsg = fmt.Sprint("game.cricket.message.open")
				} else {
					state.LastMsg = fmt.Sprint("game.cricket.message.close")
				}
				if open && remain > 0 {
					game.score(sector.Val, remain)
				} else {
					game.checkWinner()
				}
			} else {
				state.LastMsg = fmt.Sprint("game.cricket.message.hit")
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
		game.state.LastMsg = fmt.Sprint("game.message.score")
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
	game.state.LastMsg = fmt.Sprint("game.message.winner")
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	state.LastMsg = fmt.Sprint("game.message.rank")
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

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *Cricket) HoldOrNextPlayer() {
	commonHoldOrNextPlayer(game)
}

func (game *Cricket) nextDart() {
	commonNextDart(game)
}

func (game *Cricket) nextPlayer() {
	commonNextPlayer(game)
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

var gsCricketOptions = []common.GameOption{
	{"NoScore", "bool", "game.cricket.options.noscore", false},
	{"CutThroat", "bool", "game.cricket.options.cutthroat", false}}

// GsCricket GameStyle for Cricket games
var GsCricket = common.GameStyle{
	"game.cricket.name",
	"CRICKET",
	"game.cricket.rules",
	gsCricketOptions}

func newOptionCricket(opts map[string]interface{}) OptionCricket {
	o := OptionCricket{}
	gameOptionFiller(&o, gsCricketOptions, opts)

	return o
}
