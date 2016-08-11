package client

import "github.com/gocaine/go-dart/common"

//Failure is structure for failure response
type Failure struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
}

//GameResponse is container structure for GameState
type GameResponse struct {
	ID   int       `json:"id,omitempty"`
	Game GameState `json:"game"`
}

//GameState is structure for state
type GameState struct {
	State common.GameState
}
