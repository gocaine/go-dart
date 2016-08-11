package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
	"github.com/gocaine/go-dart/server/autogen"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

// Server is used to handle games
type Server struct {
	boards []string
	games  map[int]game.Game
	hubs   map[int]*GameHub
}

// NewServer Server instantiation
func NewServer() *Server {
	server := new(Server)
	server.boards = make([]string, 0)
	server.games = make(map[int]game.Game)
	server.hubs = make(map[int]*GameHub)
	return server
}

// Start prepares and start the web server
func (server *Server) Start() {
	fmt.Println("Ready to Dart !!")
	engine := gin.Default()

	apiRouter := engine.Group("api")

	// les styles de jeu possibles
	apiRouter.GET("/styles", server.getStylesHandler) // retourne la liste des styles
	// creation du jeu (POST) -  fournit le type de jeu
	apiRouter.POST("/games", server.createNewGameHandler) // retourne un id
	// board registration (POST)
	apiRouter.POST("/boards", server.registerBoardHandler) // return 202 if ok
	// registered boards list (GET)
	apiRouter.GET("/boards", server.registeredBoardsListHandler) // return registered boards
	// retourne la liste des jeux (GET) -  fournit le type de jeu
	apiRouter.GET("/games", server.listeGamesHandler) // retourne un id
	// etat du jeu (GET)
	apiRouter.GET("/games/:gameId", server.findGameByIDHandler)
	// // creation du joueur (POST) -> retourne joueur
	apiRouter.POST("/games/:gameId/players", server.addPlayerToGameHandler)
	// // etat joueur
	// r.GET("/games/{gameId}/user/{userId}", server.userHandler).Methods("GET")
	//
	// // POST : etat de la flechette
	apiRouter.POST("/darts", server.dartHandler)

	apiRouter.GET("/games/:gameId/ws", server.wsHandler)

	assetsRouter := engine.Group("/web")
	assetsRouter.StaticFS("/", &assetfs.AssetFS{Asset: autogen.Asset, AssetDir: autogen.AssetDir, AssetInfo: autogen.AssetInfo, Prefix: "webapp/dist"})

	engine.Any("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "web")
	})
	engine.Run(":8080")
}

func (server *Server) wsHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}

	log.WithFields(log.Fields{"gameID": gameID}).Info("flushing game w/ id")
	wsHandler := websocket.Handler(server.hubs[gameID].handle)
	wsHandler.ServeHTTP(c.Writer, c.Request)
}

func (server *Server) listeGamesHandler(c *gin.Context) {
	ids := make([]int, 0, len(server.games))
	for k := range server.games {
		ids = append(ids, k)
	}
	c.JSON(http.StatusOK, ids)
}

func (server *Server) registeredBoardsListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, server.boards)
}

func (server *Server) registerBoardHandler(c *gin.Context) {
	var b common.BoardRepresentation
	if c.BindJSON(&b) == nil {

		for _, board := range server.boards {
			if board == b.Name {
				c.JSON(http.StatusForbidden, gin.H{"status": "Already registered"})
				return
			}
		}

		server.boards = append(server.boards, b.Name)
		c.JSON(http.StatusAccepted, gin.H{})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content"})
	}
}

///GamesHandler
func (server *Server) createNewGameHandler(c *gin.Context) {
	var g common.GameRepresentation
	if c.BindJSON(&g) == nil {
		nextID := len(server.games) + 1

		theGame, err := gameFactory(g.Style)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
			return
		}
		server.games[nextID] = theGame
		server.hubs[nextID] = NewGameHub(theGame)

		c.JSON(http.StatusCreated, gin.H{"id": nextID, "game": theGame})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content"})
	}
}

func gameFactory(style string) (result game.Game, err error) {
	switch style {
	case common.Gs301.Code:
		result = game.NewGamex01(game.Optionx01{Score: 301, DoubleOut: false})
		return
	case common.Gs301DO.Code:
		result = game.NewGamex01(game.Optionx01{Score: 301, DoubleOut: true})
		return
	case common.Gs501.Code:
		result = game.NewGamex01(game.Optionx01{Score: 501, DoubleOut: false})
		return
	case common.Gs501DO.Code:
		result = game.NewGamex01(game.Optionx01{Score: 501, DoubleOut: true})
		return
	case common.GsHigh3.Code:
		result = game.NewGameHighest(game.OptionHighest{Rounds: 3})
		return
	case common.GsHigh5.Code:
		result = game.NewGameHighest(game.OptionHighest{Rounds: 5})
		return
	case common.GsCountup300.Code:
		result = game.NewGameCountUp(game.OptionCountUp{Target: 300})
		return
	case common.GsCountup500.Code:
		result = game.NewGameCountUp(game.OptionCountUp{Target: 500})
		return
	case common.GsCountup900.Code:
		result = game.NewGameCountUp(game.OptionCountUp{Target: 900})
		return
	case common.GsCricket.Code:
		result = game.NewGameCricket(game.OptionCricket{})
		return
	case common.GsCricketCutThroat.Code:
		result = game.NewGameCricket(game.OptionCricket{CutThroat: true})
		return
	case common.GsCricketNoScore.Code:
		result = game.NewGameCricket(game.OptionCricket{NoScore: true})
		return
	default:
		err = errors.New("game of type " + style + " is not yet supported")
		return
	}
}

func (server *Server) findGameByIDHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}
	log.WithFields(log.Fields{"gameID": gameID}).Info("flushing game w/ id")

	currentGame, ok := server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": currentGame})
}

func (server *Server) addPlayerToGameHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}

	log.Infof("flushing game w/ id {}", gameID)

	currentGame, ok := server.games[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	// TODO check that the board is not already used by another player in a different game

	var p common.PlayerRepresentation
	if c.BindJSON(&p) == nil {
		currentGame.AddPlayer(p.Board, p.Name)
		c.JSON(http.StatusCreated, "http://localhost:8080/games/"+strconv.Itoa(gameID)+"/players")
		server.hubs[gameID].refresh()
	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
}

func (server *Server) dartHandler(c *gin.Context) {

	log.Info("throwed dart")

	var d common.DartRepresentation
	if c.BindJSON(&d) == nil {

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
			server.hubs[currentGameID].refresh()
		}

	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
}

func (server *Server) getStylesHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"styles": common.GsStyles})
}
