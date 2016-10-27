package server

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
)

// Server is used to handle games
type Server struct {
	boards      map[string]*board
	games       map[int]game.Game
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
	server.games = make(map[int]game.Game)
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

	// list game flavours
	apiRouter.GET("/styles", server.getStylesHandler)

	// create a new game
	apiRouter.POST("/games", server.createNewGameHandler)

	// register a new board
	apiRouter.POST("/boards", server.registerBoardHandler)

	// returns the registered boards
	apiRouter.GET("/boards", server.registeredBoardsListHandler)

	// returns the list of games
	apiRouter.GET("/games", server.listeGamesHandler)

	// get the current state of the game
	apiRouter.GET("/games/:gameId", server.findGameByIDHandler)

	// cancel a game and free the boards
	apiRouter.DELETE("/games/:gameId", server.cancelGameHandler)

	// add player to the game
	apiRouter.POST("/games/:gameId/players", server.addPlayerToGameHandler)

	// hold while removing darts
	apiRouter.POST("/games/:gameId/hold", server.holdOrNextPlayerHandler)

	// a dart hit the board !
	apiRouter.POST("/darts", server.dartHandler)

	//get a websocket to listen for game events
	apiRouter.GET("/games/:gameId/ws", server.wsHandler)

	// serve the static files
	engine.Use(ServeStatics())
	engine.NoRoute(RerouteToIndex("/static", "/api"))

	engine.Run(":" + port)
}

///GamesHandler
func (server *Server) findGameByIDHandler(c *gin.Context) {
	server.removeEndedGame()

	if _, currentGame, ok := server.findGame(c); ok {
		c.JSON(http.StatusOK, gin.H{"game": currentGame.State()})
	}
}

func (server *Server) publishUpdate(gameID int) {
	server.hubs[gameID].refresh()
}

func (server *Server) isBoardRegistered(board string) bool {
	_, found := server.boards[board]
	return found
}

func (server *Server) removeEndedGame() {
	log.Info("removeEndedGame")
	for gameID, game := range server.games {
		// Game is over so we delete it
		if game.State().Ongoing == common.OVER {
			log.WithFields(log.Fields{"gameID": gameID}).Info("remove ended game")
			// we remove it from active Games
			for board, idGame := range server.activeGames {
				if gameID == idGame {
					delete(server.activeGames, board)
				}
			}
			hub, ok := server.hubs[gameID]
			if ok {
				hub.close()
				delete(server.hubs, gameID)
			}
			delete(server.games, gameID)
		}
	}
}
