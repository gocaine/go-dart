package common

type SectorName int
type SectorPos int

const (
	OUT SectorName = iota
	ONE SectorName = iota
	TWO SectorName = iota
	THREE SectorName = iota
	FOUR SectorName = iota
	FIVE SectorName = iota
	SIX SectorName = iota
	SEPT SectorName = iota
	EIGHT SectorName = iota
	NINE SectorName = iota
	TEN SectorName = iota
	ELEVEN SectorName = iota
	TWELVE SectorName = iota
	THIRTEEN SectorName = iota
	FOURTEEN SectorName = iota
	FIFTEEN SectorName = iota
	SIXTEEN SectorName = iota
	SEVENTEEN SectorName = iota
	EIGHTEEN SectorName = iota
	NINETEEN SectorName = iota
	TWENTY SectorName = iota
	BULLEDEYE SectorName = iota
)

const (
	_ SectorPos = iota
	SIMPLE SectorPos = iota
	DOUBLE SectorPos = iota
	TRIPLE SectorPos = iota
)

type Sector struct {
	Name SectorName
	Pos SectorPos
}

type GameState struct {

}