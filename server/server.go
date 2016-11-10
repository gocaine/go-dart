package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
	"github.com/gocaine/go-dart/i18n"
	"golang.org/x/net/websocket"
)

// Server is used to handle games
type Server struct {
	boards      map[string]*board
	hubs        map[int]*GameHub
	activeGames map[string]int
}

type board struct {
	name     string
	lastSeen time.Time
}

// NewServer Server instantiation
func NewServer() *Server {
	server := new(Server)
	server.boards = make(map[string]*board)
	server.hubs = make(map[int]*GameHub)
	server.activeGames = make(map[string]int)
	server.watchdog()
	return server
}

// Start prepares and start the web server
func (server *Server) Start(port string) {
	fmt.Println("Ready to Dart !!")
	engine := gin.Default()

	apiRouter := engine.Group("/api")

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
	// cancel a Game
	apiRouter.DELETE("/games/:gameId", server.cancelGameHandler)
	// // creation du joueur (POST) -> retourne joueur
	apiRouter.POST("/games/:gameId/players", server.addPlayerToGameHandler)
	// hold or next player (POST) - return new state
	apiRouter.POST("/games/:gameId/hold", server.holdOrNextPlayerHandler)
	// // etat joueur
	// r.GET("/games/{gameId}/user/{userId}", server.userHandler).Methods("GET")
	//
	// // POST : etat de la flechette
	apiRouter.POST("/darts", server.dartHandler)

	apiRouter.GET("/games/:gameId/ws", server.wsHandler)

	engine.Use(ServeStatics())
	engine.NoRoute(RerouteToIndex("/static", "/api"))

	engine.Run(":" + port)
}

func (server *Server) wsHandler(c *gin.Context) {

	server.removeEndedGame()

	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}

	_, ok := server.hubs[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	ctx := createContext(c)
	log.WithFields(log.Fields{"gameID": gameID}).Info("flushing game w/ id")
	wsHandler := websocket.Handler(server.hubs[gameID].handle(ctx.Locale))
	wsHandler.ServeHTTP(c.Writer, c.Request)
}

func (server *Server) listeGamesHandler(c *gin.Context) {
	server.removeEndedGame()
	ids := make([]int, 0, len(server.hubs))
	for k := range server.hubs {
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
					gc := common.GameContext{}
					gc.Locale = i18n.GetLocale("eng")
					for gameID, game := range server.hubs {
						if !game.BoardHasLeft(gc, board.name) {
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

///GamesHandler

func (server *Server) cancelGameHandler(c *gin.Context) {
	server.removeEndedGame()

	if gameID, currentGame, ok := server.findGame(c); ok {
		currentGame.State().Ongoing = common.OVER
		server.publishUpdate(gameID)
	}
}

func (server *Server) createNewGameHandler(c *gin.Context) {
	server.removeEndedGame()
	var g common.NewGameRepresentation
	if c.BindJSON(&g) == nil {
		nextID := len(server.hubs) + 1

		gc := createContext(c)

		theGame, err := gameFactory(gc, g)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
			return
		}
		server.hubs[nextID] = NewGameHub(theGame)

		c.JSON(http.StatusCreated, gin.H{"id": nextID, "game": theGame})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content"})
	}
}

func gameFactory(ctx common.GameContext, g common.NewGameRepresentation) (result game.Game, err error) {

	switch g.Style {
	case game.GsX01.Code:
		result, err = game.NewGamex01(ctx, g.Options)
	case game.GsCountUp.Code:
		result, err = game.NewGameCountUp(ctx, g.Options)
	case game.GsHighest.Code:
		result, err = game.NewGameHighest(ctx, g.Options)
	case game.GsCricket.Code:
		result, err = game.NewGameCricket(ctx, g.Options)
	default:
		err = errors.New("game of type " + g.Style + " is not yet supported")
	}
	return
}

func (server *Server) findGameByIDHandler(c *gin.Context) {
	server.removeEndedGame()

	if _, currentGame, ok := server.findGame(c); ok {
		c.JSON(http.StatusOK, gin.H{"game": currentGame.State()})
	}
}

func (server *Server) addPlayerToGameHandler(c *gin.Context) {
	server.removeEndedGame()

	if gameID, currentGame, ok := server.findGame(c); ok {
		var p common.PlayerRepresentation
		if c.BindJSON(&p) == nil {

			if !server.isBoardRegistered(p.Board) {
				c.JSON(http.StatusNotFound, gin.H{"status": "board not found"})
				return
			}

			activeGameID, ok := server.activeGames[p.Board]
			if ok && activeGameID != gameID {
				// A active game, different from this one, has been found
				c.JSON(http.StatusForbidden, gin.H{"status": "board is already busy"})
				return
			}

			gc := createContext(c)
			currentGame.AddPlayer(gc, p.Board, p.Name)
			server.activeGames[p.Board] = gameID
			c.JSON(http.StatusCreated, nil)
			server.hubs[gameID].refresh()
		} else {
			c.JSON(http.StatusBadRequest, nil)
		}
	}

}

func (server *Server) holdOrNextPlayerHandler(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}

	currentGame, ok := server.hubs[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	gc := createContext(c)
	currentGame.HoldOrNextPlayer(gc)
	state := currentGame.State()
	c.JSON(http.StatusOK, gin.H{"state": state})
	server.hubs[gameID].refresh()
}

func (server *Server) dartHandler(c *gin.Context) {

	var d common.DartRepresentation
	if c.BindJSON(&d) == nil {

		log.Infof("revieved a dart %v", d)

		var currentGame game.Game
		var currentGameID int
		for gameID, game := range server.hubs {
			if game.State().Players[game.State().CurrentPlayer].Board == d.Board {
				currentGame = game
				currentGameID = gameID
			}
		}

		if currentGame == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		gc := createContext(c)

		state, err := currentGame.HandleDart(gc, common.Sector{Val: d.Sector, Pos: d.Multiplier})
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

func (server *Server) findGame(c *gin.Context) (gameID int, currentGame game.Game, ok bool) {

	gameID, err := strconv.Atoi(c.Param("gameId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content", "error": err.Error()})
		return
	}
	log.WithFields(log.Fields{"gameID": gameID}).Info("flushing game w/ id")

	currentGame, ok = server.hubs[gameID]
	if !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	return
}

func (server *Server) publishUpdate(gameID int) {
	server.hubs[gameID].refresh()
}

func (server *Server) getStylesHandler(c *gin.Context) {

	gc := createContext(c)

	c.JSON(http.StatusOK, gin.H{"styles": game.Flavors(gc)})
}

func (server *Server) isBoardRegistered(board string) bool {
	_, found := server.boards[board]
	return found
}

func (server *Server) removeEndedGame() {
	log.Info("removeEndedGame")
	for gameID, hub := range server.hubs {
		// Game is over so we delete it
		if hub.State().Ongoing == common.OVER {
			log.WithFields(log.Fields{"gameID": gameID}).Info("removeEndedGame")
			// we remove it from active Games
			for board, idGame := range server.activeGames {
				if gameID == idGame {
					delete(server.activeGames, board)
				}
			}
			hub.close()
			delete(server.hubs, gameID)
		}
	}
}

func createContext(c *gin.Context) common.GameContext {
	gc := common.GameContext{}
	gc.Locale = i18n.GetLocale(c.Request.Header.Get("Accept-Language"))
	return gc

}
