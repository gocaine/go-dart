package client

import "github.com/gocaine/go-dart/common"

//Failure is structure for failure response
type Failure struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
}

//GameResponse is container structure for GameState
type GameResponse struct {
	Id   int       `json:"id,omitempty"`
	Game GameState `json:"game"`
}

//GameResponse is structure for state
type GameState struct {
	State common.GameState
}
