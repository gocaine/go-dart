package server

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func (server *Server) wsHandler(c *gin.Context) {

	server.removeEndedGame()

	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}

	_, ok := server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	log.WithFields(log.Fields{"gameID": gameID}).Info("flushing game w/ id")
	wsHandler := websocket.Handler(server.hubs[gameID].handle)
	wsHandler.ServeHTTP(c.Writer, c.Request)
}
