package i18n

import "github.com/olebedev/config"

var baseConfig *config.Config

var yamlFiles = make(map[string]string)

var baseYaml = `
game:
  message:
    rank: Player %d end at rank #%d
    disconnect: Board %s has been disconnected
    score:
      "Scored : %d"
    winner:
      "Winner : %s"
    player:
      exists: Player name is already in use
      notadded: Player cannot be added
      next: Next Player
  error:
    onhold: Game is on hold and not ready to handle darts
    notstarted: Game is not started or is ended
    cantstart: Game cannot start
    sector:
      invalid: Sector is not a valid one
  countup:
    name: CountUp
    display: Count-Up %d
    error:
      target: Target should be at least 61
    rules: All players start with 0 points and attempt to reach the given target (300 / 500 / ...).
    options:
      target: The score to reach
  cricket:
    name: Cricket
    display:
      cricket: Cricket
      cutthroat: Cut-Throat Cricket
      noscore: No Score Cricket
    message:
      open:
        "Opened : %s"
      close:
        "Closed : %s"
      hit:
        "Hit : %d x %s"
    error:
      incompatible: CutThroat and NoScore options are not compatible
    options:
      noscore: If set to true, no point is scored, the winner is the first player to close all sectors
      cutthroat: If set to true, when a player hit a sector for the 4th time or more, the points go to the players who havent close the sector. In the end, the winner is the first to close every sector with the smallest score
    rules: The main purpose is to open (or close) all the sectors. The sectors are 15, 16, 17, 18, 19, 20 and bull's eye. To open a sector a player has to hit it 3 times (a Triple counts for 3 hits, a Double for 2). When a sector is open for a player, he can score in it (the points are the real value). When all players have open a given sector it is close, and no more point are scored in it. The winner is the first player to both have open all the sectors and the highest score
  highest:
    name: Highest
    display: "%d visits HighScore"
    error:
      rounds: Rounds should be at least 1
    options:
      rounds: The number of visits each player play
    rules: All players throw the same number of darts (3 per rounds) then the player with the highest score wins
  x01:
    name:
      "X01 : 301, 501,..."
    display:
      x01: "%d"
      doubleout: "%d Double-out"
    message:
      doubleout: You should end with a double
      overscore: You went beyond the target dude !
    error:
      score: Score should be at least 61
    options:
      score: The score from which to reach 0
      doubleout: If set to true, the players have to end with a double (and so reach 0)
    rules:
      All players start with the same points (301 / 501 / ...) and attempt to reach zero. If a player scores more than the total required to reach zero, the player "busts" and the score returns to the score that was existing at the start of the turn.
`

func init() {
	var err error
	baseConfig, err = config.ParseYaml(baseYaml)
	if err != nil {
		panic(err)
	}
}
