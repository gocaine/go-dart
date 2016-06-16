package server

import "go-dart/common"

type Game interface {
	Start() error
	AddPlayer(name string) error
	HandleDart(sector common.Sector) (*common.GameState, error)
}
