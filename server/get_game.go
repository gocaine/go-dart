package server

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/game"
)

func (server *Server) findGame(c *gin.Context) (gameID int, currentGame game.Game, ok bool) {

	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}
	log.WithFields(log.Fields{"gameID": gameID}).Info("flushing game w/ id")

	currentGame, ok = server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	return
}
