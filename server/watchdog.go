package server

import (
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/gocaine/go-dart/common"
)

func (server *Server) watchdog() {
	log.Debugf("Starting watchdog (timeout: %v)", common.HealthCheckTimeout)
	ticker := time.NewTicker(common.HealthCheckTimeout)
	go func() {
		for {
			<-ticker.C
			deadline := time.Now().Add(-common.HealthCheckTimeout)
			log.Debugf("healthcheck deadline is %v", deadline)

			for name, board := range server.boards {
				log.Debugf("board %s was last seen %v", board.name, board.lastSeen)
				if board.lastSeen.Before(deadline) {
					log.Warnf("board %s seems to be dead", board.name)
					delete(server.boards, name)
					hasEndedGame := false
					for gameID, game := range server.games {
						if !game.BoardHasLeft(board.name) {
							server.publishUpdate(gameID)
							hasEndedGame = true
						}
					}

					if hasEndedGame {
						server.removeEndedGame()
					}
				}
			}
		}
	}()
}
