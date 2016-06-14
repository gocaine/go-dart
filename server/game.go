package server

import "go-dart/common"

type Game interface {
	Start()
	AddPlayer(name string)
	HandleDart(sector common.Sector) common.GameState
}
