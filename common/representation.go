package common

// GameRepresentation Json data representation for a Game
type GameRepresentation struct {
	Style string `json:"style"`
}

// PlayerRepresentation Json data representation for a Player
type PlayerRepresentation struct {
	Name  string `json:"name"`
	Board string `json:"board"`
}

// DartRepresentation Json data representation for a Dart
type DartRepresentation struct {
	Sector     int    `json:"sector"`
	Multiplier int    `json:"multiplier"`
	Board      string `json:"board"`
}

// BoardRepresentation Json data representation for a Board
type BoardRepresentation struct {
	Name string `json:"name"`
}
