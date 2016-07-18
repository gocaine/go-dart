# go-dart

[![Build Status](https://travis-ci.org/gocaine/go-dart.svg?branch=master)](https://travis-ci.org/gocaine/go-dart)
[![Coverage Status](https://coveralls.io/repos/github/gocaine/go-dart/badge.svg?branch=master)](https://coveralls.io/github/gocaine/go-dart?branch=master)

# Authors

- Guillaume GERBAUD
- Mathieu POUSSE
- Erwann THEBAULT
- Jeremie HUCHET
- Maximilien RICHER

# Bundle resources

First generate resources

```
cd webapp
#npm install
docker run --rm -v $PWD:/data ggerbaud/node-bower-grunt:5 npm install
# bower install
docker run --rm -v $PWD:/data ggerbaud/node-bower-grunt:5 bower install
# grunt build
docker run --rm -v $PWD:/data ggerbaud/node-bower-grunt:5 grunt build
```

Then generate bundle

`esc -o server/statics.go -pkg="server" -prefix="webapp/dist" webapp/dist`

(esc should be in $GOPATH/bin)


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

