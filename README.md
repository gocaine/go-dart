# go-dart

[![Build Status](https://travis-ci.org/Zenika/go-dart.svg?branch=master)](https://travis-ci.org/Zenika/go-dart)

# Contributors

- Guillaume GERBAUD
- Jeremie HUCHET
- Mathieu POUSSE
- Maximilien RICHER
- Erwann THEBAULT

# API

- Create a game
  + `POST "/games"`
  + return game ID
- Get the current game state
  + `GET "/games/{id}"`
  + return a GameState
- Create player
  + `POST "/games/{id}/players"`
  + return User ID
- Player state
  + `GET "/games/{id}/players/{id}"`
- Dart state
  + `POST "/games/{id}/dart"`

# Scenario

Create a new game

    curl -X POST -d '{"style": "301"}' http://localhost:8080/games

Add players

    curl -X POST -d '{"name": "player 1"}' http://localhost:8080/games/1/players
    curl -X POST -d '{"name": "player 2"}' http://localhost:8080/games/1/players

Throw darts

    curl -X POST -d '{"sector": 20, "multiplier": 1}' http://localhost:8080/games/1/darts
    curl -X POST -d '{"sector": 20, "multiplier": 2}' http://localhost:8080/games/1/darts
    curl -X POST -d '{"sector": 19, "multiplier": 1}' http://localhost:8080/games/1/darts
