package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) listeGamesHandler(c *gin.Context) {
	server.removeEndedGame()
	ids := make([]int, 0, len(server.games))
	for k := range server.games {
		ids = append(ids, k)
	}
	c.JSON(http.StatusOK, ids)
}

func (server *Server) registeredBoardsListHandler(c *gin.Context) {
	boards := make([]string, len(server.boards))
	i := 0
	for k := range server.boards {
		boards[i] = k
		i++
	}
	c.JSON(http.StatusOK, boards)
}
