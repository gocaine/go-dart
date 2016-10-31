package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) holdOrNextPlayerHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}

	currentGame, ok := server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	currentGame.HoldOrNextPlayer()
	state := currentGame.State()
	c.JSON(http.StatusOK, gin.H{"state": state})
	server.hubs[gameID].refresh()
}
