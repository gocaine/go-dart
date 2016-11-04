package game

import (
	"errors"
	"sort"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/i18n"
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
func NewGameCricket(ctx common.GameContext, opts map[string]interface{}) (g *Cricket, err error) {
	opt := newOptionCricket(opts)
	if opt.CutThroat && opt.NoScore {
		err = errors.New(i18n.Translation("game.cricket.error.incompatible", ctx.Locale))
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
	g.DisplayStyle = i18n.Translation(dStyle, ctx.Locale)
	g.memory = make(map[string]int)

	return
}

// AddPlayer add a new player to the game
func (game *Cricket) AddPlayer(ctx common.GameContext, board string, name string) (error error) {

	error = commonAddPlayer(ctx, game, board, name)
	if error == nil {
		game.state.Players[len(game.state.Players)-1].Histo = make(map[string]int)
	}
	log.WithFields(log.Fields{"name": name, "player": game.state.Players[len(game.state.Players)-1]}).Info("AddPlayer")
	return
}

// Start start the game, Darts will be handled
func (game *Cricket) Start(ctx common.GameContext) (error error) {

	error = commonStart(ctx, game)
	if error == nil {
		for _, key := range sectors {
			game.memory[key] = len(game.state.Players)
		}
		log.WithFields(log.Fields{"memory": game.memory}).Info("Start")
	}
	return
}

// HandleDart the implementation has to handle the Dart regarding the current player, the cricket rules, and the context. Return a GameState
func (game *Cricket) HandleDart(ctx common.GameContext, sector common.Sector) (result *common.GameState, error error) {

	error = commonHandleDartChecks(ctx, game, sector)
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
				game.score(ctx, sector.Val, sector.Pos)
			} else {
				game.NextDart(ctx)
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
					ctx.MessageHandler("game.cricket.message.open", sVal)
				} else {
					ctx.MessageHandler("game.cricket.message.close", sVal)
				}
				if open && remain > 0 {
					game.score(ctx, sector.Val, remain)
				} else {
					game.checkWinner(ctx)
				}
			} else {
				ctx.MessageHandler("game.cricket.message.hit", sector.Pos, sVal)
				game.NextDart(ctx)
			}
		}
	} else {
		game.NextDart(ctx)
	}
	result = state
	return
}

func (game *Cricket) score(ctx common.GameContext, val, pos int) {
	log.WithFields(log.Fields{"sector": val, "number": pos}).Info("score")
	if game.noScore {
		// no score at all
	} else {
		points := val * pos
		ctx.MessageHandler("game.message.score", points)
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

	game.checkWinner(ctx)
}

func (game *Cricket) checkWinner(ctx common.GameContext) {
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
			game.winner(ctx)
		} else {
			if game.cutThroat {
				if lowest(game.state.Players, game.state.CurrentPlayer) {
					game.winner(ctx)
				} else {
					game.NextDart(ctx)
				}
			} else {
				if highest(game.state.Players, game.state.CurrentPlayer) {
					game.winner(ctx)
				} else {
					game.NextDart(ctx)
				}
			}
		}
	} else {
		game.NextDart(ctx)
	}
}

func (game *Cricket) winner(ctx common.GameContext) {
	log.Info("winner")
	state := game.state
	ctx.MessageHandler("game.message.winner", state.Players[state.CurrentPlayer].Name)
	state.Players[state.CurrentPlayer].Rank = game.rank + 1
	ctx.MessageHandler("game.message.rank", state.CurrentPlayer, game.rank+1)
	game.rank++
	if game.rank >= len(state.Players)-1 {
		game.state.Ongoing = common.OVER
		sort.Sort(common.ByRank(state.Players))
		if len(state.Players) > 1 {
			state.Players[len(state.Players)-1].Rank = game.rank + 1
		}
	} else {
		game.NextPlayer(ctx)
	}
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING
func (game *Cricket) HoldOrNextPlayer(ctx common.GameContext) {
	commonHoldOrNextPlayer(ctx, game)
}

// NextDart is called after each dart when the same palyer play again
func (game *Cricket) NextDart(ctx common.GameContext) {
	commonNextDart(ctx, game)
}

// NextPlayer is called when the current player end his visit
func (game *Cricket) NextPlayer(ctx common.GameContext) {
	commonNextPlayer(ctx, game)
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
