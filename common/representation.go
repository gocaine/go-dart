package common

type GameRepresentation struct {
	Style string `json:"style"`
}

type PlayerRepresentation struct {
	Name string `json:"name"`
}

type DartRepresentation struct {
	Sector     int `json:"sector"`
	Multiplier int `json:"multiplier"`
}
