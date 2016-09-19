package game

import (
	"fmt"
	"sort"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/pkg/errors"
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
		err = errors.New("CutThroat and NoScore options are not compatible")
	}
	g = new(Cricket)
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
	{"NoScore", "bool", "If set to true, no point is scored, the winner is the first player to close all sectors", false},
	{"CutThroat", "bool", "If set to true, when a player hit a sector for the 4th time or more, the points go to the players who havent close the sector. " +
		"In the end, the winner is the first to close every sector with the smallest score", false}}

// GsCricket GameStyle for Cricket games
var GsCricket = common.NewGameStyle{
	"Cricket",
	"CRICKET",
	"The main purpose is to open (or close) all the sectors. The sectors are 15, 16, 17, 18, 19, 20 and bull's eye." +
		" To open a sector a player has to hit it 3 times (a Triple counts for 3 hits, a Double for 2). When a sector is open for a player, he can score in it (the points are the real value). " +
		"When all players have open a given sector it is close, and no more point are scored in it. " +
		"The winner is the first player to both have open all the sectors and the highest score",
	gsCricketOptions}

func newOptionCricket(opts map[string]interface{}) OptionCricket {
	o := OptionCricket{}
	gameOptionFiller(&o, gsCricketOptions, opts)

	return o
}
