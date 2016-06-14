package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
}

func NewServer() *Server {
	server := new(Server)
	return server
}

func (*Server) Start() {
	fmt.Println("Ready to Dart !!")

	r := mux.NewRouter()
	// creation du jeu (POST) -  fournit le type de jeu
	r.HandleFunc("/games", GamesHandler).Methods("POST") // retourne un id
	// etat du jeu (GET)
	r.HandleFunc("/games/{gameId}", GameHandler).Methods("GET")
	// creation du joueur (POST) -> retourne joueur
	r.HandleFunc("/games/{gameId}/user", UsersHandler).Methods("POST")
	// etat joueur
	r.HandleFunc("/games/{gameId}/user/{userId}", UserHandler).Methods("GET")

	// POST : etat de la flechette
	r.HandleFunc("/games/{gameId}/dart", DartHandler).Methods("POST")

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}

type gameRepresentation struct {
	Style string `json:"style"`
}

///GamesHandler
func GamesHandler(writer http.ResponseWriter, request *http.Request) {
	var g gameRepresentation
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&g)
	fmt.Fprintf(writer, "yeah %s ! ", g.Style)
}

func GameHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameId"]

	result, _ := json.Marshal(gameID)

	fmt.Fprint(writer, "gameID "+string(result))
}

func UsersHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameID"]
	fmt.Fprint(writer, "gameID "+gameID)
}

func UserHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameId"]
	userID := vars["userId"]
	fmt.Fprint(writer, "gameID "+gameID+" userId"+userID)
}

func DartHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	gameID := vars["gameId"]
	fmt.Fprint(writer, "gameID "+gameID+" dart")
}
