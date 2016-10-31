package server

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
)

func (server *Server) registerBoardHandler(c *gin.Context) {
	server.removeEndedGame()
	var representation common.BoardRepresentation
	if c.BindJSON(&representation) == nil {
		existing, found := server.boards[representation.Name]
		if found {
			existing.lastSeen = time.Now()
			log.Debugf("pong %s (last seen is now %v)", representation.Name, existing.lastSeen)
			c.JSON(http.StatusOK, gin.H{"status": "Pong"})
		} else {
			log.Infof("new board has been registered: %s", representation.Name)
			server.boards[representation.Name] = &board{name: representation.Name, lastSeen: time.Now()}
			c.JSON(http.StatusAccepted, gin.H{})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content"})
	}
}
