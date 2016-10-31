package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
)

func (server *Server) dartHandler(c *gin.Context) {

	var d common.DartRepresentation
	if c.BindJSON(&d) == nil {

		log.Infof("received a dart %v", d)

		var currentGame game.Game
		var currentGameID int
		for gameID, game := range server.games {
			if game.State().Players[game.State().CurrentPlayer].Board == d.Board {
				currentGame = game
				currentGameID = gameID
			}
		}

		if currentGame == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		state, err := currentGame.HandleDart(common.Sector{Val: d.Sector, Pos: d.Multiplier})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"state": state})
			server.publishUpdate(currentGameID)
		}

	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
}
