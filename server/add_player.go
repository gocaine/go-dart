package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
)

func (server *Server) addPlayerToGameHandler(c *gin.Context) {
	server.removeEndedGame()

	if gameID, currentGame, ok := server.findGame(c); ok {
		var p common.PlayerRepresentation
		if c.BindJSON(&p) == nil {

			if !server.isBoardRegistered(p.Board) {
				c.JSON(http.StatusNotFound, gin.H{"status": "board not found"})
				return
			}

			activeGameID, ok := server.activeGames[p.Board]
			if ok && activeGameID != gameID {
				// A active game, different from this one, has been found
				c.JSON(http.StatusForbidden, gin.H{"status": "board is already busy"})
				return
			}

			currentGame.AddPlayer(p.Board, p.Name)
			server.activeGames[p.Board] = gameID
			c.JSON(http.StatusCreated, nil)
			server.hubs[gameID].refresh()
		} else {
			c.JSON(http.StatusBadRequest, nil)
		}
	}

}
