package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func start() {
	fmt.Println("Ready to Dart !!")

	r := mux.NewRouter()
	// creation du jeu (POST)
	r.HandleFunc("/games", GamesHandler) // retourne un id
	// etat du jeu (GET)
	r.HandleFunc("/games/{id}", GamesHandler)
	// creation du joueur (POST) -> retourne joueur
	r.HandleFunc("/games/{id}/user", GamesHandler)
	// etat joueur
	r.HandleFunc("/games/{id}/user/{id}", GamesHandler)

	// POST : etat de la flechette
	r.HandleFunc("/games/{id}/dart", GamesHandler)

	r.HandleFunc("/games", GamesHandler)
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}

func GamesHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	name := vars["name"]
	fmt.Fprint(writer, "yeah "+name)
}
