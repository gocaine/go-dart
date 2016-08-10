package common

type GameRepresentation struct {
	Style string `json:"style"`
	Board string `json:"board"`
}

type PlayerRepresentation struct {
	Name string `json:"name"`
}

type DartRepresentation struct {
	Sector     int    `json:"sector"`
	Multiplier int    `json:"multiplier"`
	Board      string `json:"board"`
}
