package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
)

func (server *Server) getStylesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"styles": [...]common.GameStyle{game.GsX01, game.GsCountUp, game.GsHighest, game.GsCricket}})
}
