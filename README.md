# go-dart

[![Build Status](https://travis-ci.org/gocaine/go-dart.svg?branch=master)](https://travis-ci.org/gocaine/go-dart)
[![Coverage Status](https://coveralls.io/repos/github/gocaine/go-dart/badge.svg?branch=master)](https://coveralls.io/github/gocaine/go-dart?branch=master)

# Authors

- Guillaume GERBAUD
- Mathieu POUSSE
- Erwann THEBAULT

# Build

First of all, ensure you have checked out the project your `${GOPATH}/github.com/gocaine/go-dart`.

Simply run the command `make binary`

As the build process relies on docker, please first configure it on your local.

If you don't want to use docker, simply run `make dev binary` rely on your local version of `go` and `npm`.

## Contributors

- Jeremie HUCHET
- Maximilien RICHER

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

