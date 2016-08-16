package server

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/game"
	"golang.org/x/net/websocket"
)

// GameHub handle websocket connections for a Game
type GameHub struct {
	clients []*websocket.Conn
	output  chan bool
	game    game.Game
}

// NewGameHub is GameHub constructor
func NewGameHub(game game.Game) *GameHub {
	hub := GameHub{game: game, clients: make([]*websocket.Conn, 0)}
	return &hub
}

func (gh *GameHub) handle(connection *websocket.Conn) {
	log.Infof("new ws connection for this user")
	gh.clients = append(gh.clients, connection)
	status, _ := json.Marshal(gh.game.State())
	connection.Write([]byte(status))
	// lock until the end of the world
	connection.Read(make([]byte, 0))
}

func (gh *GameHub) refresh() {
	status, err := json.Marshal(gh.game.State())
	if err != nil {
		log.Info("cannot serialize status")
	}
	statusAsBytes := []byte(status)
	for _, client := range gh.clients {
		log.Info("sending status")
		_, err := client.Write(statusAsBytes)
		if err != nil {
			log.Infof("error writing %v", err)
		}
	}
}

func (gh *GameHub) close() {
	log.Infof("close all websocket connections")
	for _, client := range gh.clients {
		client.Close()
	}
}
