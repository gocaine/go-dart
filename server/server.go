package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	// creation du jeu (POST) -  fournit le type de jeu
	r.HandleFunc("/games", server.gamesHandler).Methods("POST") // retourne un id
	// etat du jeu (GET)
	r.HandleFunc("/games/{gameId}", server.gameHandler).Methods("GET")
	// creation du joueur (POST) -> retourne joueur
	r.HandleFunc("/games/{gameId}/user", server.usersHandler).Methods("POST")
	// etat joueur
	r.HandleFunc("/games/{gameId}/user/{userId}", server.userHandler).Methods("GET")

	// POST : etat de la flechette
	r.HandleFunc("/games/{gameId}/dart", server.dartHandler).Methods("POST")

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}

type gameRepresentation struct {
	Style string `json:"style"`
}

///GamesHandler
func (server *Server) gamesHandler(writer http.ResponseWriter, request *http.Request) {
	var g gameRepresentation
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&g)
	nextID := len(server.games) + 1

	theGame, err := gameFactory(g.Style)

	if err != nil {
		fmt.Fprintf(writer, "go fuck yourself %s ! ", g.Style)
	}
	server.games[nextID] = theGame
	fmt.Fprintf(writer, "%d yeah %+v ! ", nextID, server.games[nextID])
}

func gameFactory(style string) (result Game, err error) {
	switch style {
	case "301":
		result = NewGame(301)
		return
	default:
		err = errors.New("prout")
		return
	}
}

func (server *Server) gameHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameId"]

	result, _ := json.Marshal(gameID)

	fmt.Fprint(writer, "gameID "+string(result))
}

func (server *Server) usersHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameID"]
	fmt.Fprint(writer, "gameID "+gameID)
}

func (server *Server) userHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameId"]
	userID := vars["userId"]
	fmt.Fprint(writer, "gameID "+gameID+" userId"+userID)
}

func (server *Server) dartHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameId"]
	fmt.Fprint(writer, "gameID "+gameID+" dart")
}
