package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
	"github.com/gocaine/go-dart/server/autogen"

	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type Server struct {
	games map[int]game.Game
	hubs  map[int]*GameHub
}

func NewServer() *Server {
	server := new(Server)
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
	apiRouter.POST("/games/:gameId/darts", server.dartHandler)

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
	case common.GS_301.Code:
		result = game.NewGamex01(game.Optionx01{Score: 301, DoubleOut: false})
		return
	case common.GS_301_DO.Code:
		result = game.NewGamex01(game.Optionx01{Score: 301, DoubleOut: true})
		return
	case common.GS_501.Code:
		result = game.NewGamex01(game.Optionx01{Score: 501, DoubleOut: false})
		return
	case common.GS_501_DO.Code:
		result = game.NewGamex01(game.Optionx01{Score: 501, DoubleOut: true})
		return
	case common.GS_HIGH_3.Code:
		result = game.NewGameHighest(game.OptionHighest{Rounds: 3})
		return
	case common.GS_HIGH_5.Code:
		result = game.NewGameHighest(game.OptionHighest{Rounds: 5})
		return
	case common.GS_COUNTUP_300.Code:
		result = game.NewGameCountUp(game.OptionCountUp{Target: 300})
		return
	case common.GS_COUNTUP_500.Code:
		result = game.NewGameCountUp(game.OptionCountUp{Target: 500})
		return
	case common.GS_COUNTUP_900.Code:
		result = game.NewGameCountUp(game.OptionCountUp{Target: 900})
		return
	case common.GS_CRICKET.Code:
		result = game.NewGameCricket(game.OptionCricket{})
		return
	case common.GS_CRICKET_CUTTHROAT.Code:
		result = game.NewGameCricket(game.OptionCricket{CutThroat: true})
		return
	case common.GS_CRICKET_NOSCORE.Code:
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

	var p common.PlayerRepresentation
	if c.BindJSON(&p) == nil {
		currentGame.AddPlayer(p.Name)
		c.JSON(http.StatusCreated, "http://localhost:8080/games/"+strconv.Itoa(gameID)+"/players")
		server.hubs[gameID].refresh()
	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
}

func (server *Server) dartHandler(c *gin.Context) {
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

	var d common.DartRepresentation
	if c.BindJSON(&d) == nil {
		state, err := currentGame.HandleDart(common.Sector{Val: d.Sector, Pos: d.Multiplier})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"state": state})
			server.hubs[gameID].refresh()
		}

	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
}

func (server *Server) getStylesHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"styles": common.GS_STYLES})
}
