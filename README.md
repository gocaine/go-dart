# go-dart

# Contributors

- Guillaume GERBAUD
- Mathieu POUSSE
- Maximilien RICHER
- Jeremie HUCHET

# API

- Create a game
  + `POST "/games"`
  + return game ID
- Get the current game state
  + `GET "/games/{id}"
  + return a GameState
- Create player
  + `POST "/games/{id}/user"`
  + return User ID 
- Player state 
  + `GET "/games/{id}/user/{id}"`
- Dart state
  + `POST "/games/{id}/dart"`
