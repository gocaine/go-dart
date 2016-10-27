package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
)

func (server *Server) cancelGameHandler(c *gin.Context) {
	server.removeEndedGame()

	if gameID, currentGame, ok := server.findGame(c); ok {
		currentGame.State().Ongoing = common.OVER
		server.publishUpdate(gameID)
	}
}
