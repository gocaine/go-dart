package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go-dart/common"
)

type Server struct {
	games map[int]Game
}

func NewServer() *Server {
	server := new(Server)
	server.games = make(map[int]Game)
	return server
}

func (server *Server) Start() {
	fmt.Println("Ready to Dart !!")
	r := gin.Default()

	// creation du jeu (POST) -  fournit le type de jeu
	r.POST("/games", server.createNewGameHandler) // retourne un id
	// etat du jeu (GET)
	r.GET("/games/:gameId", server.findGameByIdHandler)
	// // creation du joueur (POST) -> retourne joueur
	r.POST("/games/:gameId/players", server.addPlayerToGameHandler)
	// // etat joueur
	// r.GET("/games/{gameId}/user/{userId}", server.userHandler).Methods("GET")
	//
	// // POST : etat de la flechette
	r.POST("/games/:gameId/darts", server.dartHandler)

	r.Run(":8080")
}

type gameRepresentation struct {
	Style string `json:"style"`
}

///GamesHandler
func (server *Server) createNewGameHandler(c *gin.Context) {
	var g gameRepresentation
	if c.BindJSON(&g) == nil {
		nextID := len(server.games) + 1

		theGame, err := gameFactory(g.Style)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"stats": "illegal content"})
			return
		}
		server.games[nextID] = theGame

		c.JSON(http.StatusOK, gin.H{"id": nextID, "game": theGame})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"stats": "illegal content"})
	}
}

func gameFactory(style string) (result Game, err error) {
	switch style {
	case "301":
		result = NewGamex01(Optionx01{Score: 301})
		return
	default:
		err = errors.New("prout")
		return
	}
}

func (server *Server) findGameByIdHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"stats": "illegal content"})
		return
	}
	log.Infof("flushing game w/ id {}", gameID)

	currentGame, ok := server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": currentGame})
}

type playerRepresentation struct {
	Name string `json:"name"`
}

func (server *Server) addPlayerToGameHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"stats": "illegal content"})
		return
	}

	log.Infof("flushing game w/ id {}", gameID)

	currentGame, ok := server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	var p playerRepresentation
	if c.BindJSON(&p) == nil {
		currentGame.AddPlayer(p.Name)
		c.JSON(http.StatusCreated, "http://localhost:8080/games/"+strconv.Itoa(gameID)+"/players")
	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
}

type dartRepresentation struct {
	Sector     int `json:"sector"`
	Multiplier int `json:"multiplier"`
}

func (server *Server) dartHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"stats": "illegal content"})
		return
	}

	log.Infof("flushing game w/ id {}", gameID)

	currentGame, ok := server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	var d dartRepresentation
	if c.BindJSON(&d) == nil {
		currentGame.HandleDart(common.Sector{Val: d.Sector, Pos: d.Multiplier})
		c.JSON(http.StatusOK, currentGame)
	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
}
