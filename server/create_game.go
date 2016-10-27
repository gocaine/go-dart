package server

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
	"net/http"
)

func (server *Server) createNewGameHandler(c *gin.Context) {
	server.removeEndedGame()
	var g common.NewGameRepresentation
	if c.BindJSON(&g) == nil {
		nextID := len(server.games) + 1

		theGame, err := gameFactory(g)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
			return
		}
		server.games[nextID] = theGame
		server.hubs[nextID] = NewGameHub(theGame)

		c.JSON(http.StatusCreated, gin.H{"id": nextID, "game": theGame})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content"})
	}
}

func gameFactory(g common.NewGameRepresentation) (result game.Game, err error) {

	switch g.Style {
	case game.GsX01.Code:
		result, err = game.NewGamex01(g.Options)
	case game.GsCountUp.Code:
		result, err = game.NewGameCountUp(g.Options)
	case game.GsHighest.Code:
		result, err = game.NewGameHighest(g.Options)
	case game.GsCricket.Code:
		result, err = game.NewGameCricket(g.Options)
	default:
		err = errors.New("game of type " + g.Style + " is not yet supported")
	}
	return
}
