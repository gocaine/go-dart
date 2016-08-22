
const baseUrl = '/web'

const routes = {
  // home page
  home: baseUrl + '/',
  // list game flavors
  newGame: baseUrl + '/newGame',
  // list existing games
  listGames: baseUrl + '/listGames',
  // join game
  joinGameById: baseUrl + '/games/:gameId',
}

export default routes;