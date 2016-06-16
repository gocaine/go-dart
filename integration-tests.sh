#!/bin/bash

curl -X POST -d '{"style": "701"}' http://localhost:8080/games
curl -X POST -d '{"style": "301"}' http://localhost:8080/games

curl -X POST -d '{"name": "player 1"}' http://localhost:8080/games/1/players
curl -X POST -d '{"name": "player 2"}' http://localhost:8080/games/1/players

curl -X POST -d '{"sector": 20, "multiplier": 1}' http://localhost:8080/games/1/darts
curl -X POST -d '{"sector": 20, "multiplier": 1}' http://localhost:8080/games/1/darts
curl -X POST -d '{"sector": 20, "multiplier": 1}' http://localhost:8080/games/1/darts

curl -X POST -d '{"sector": 20, "multiplier": 1}' http://localhost:8080/games/1/darts
curl -X POST -d '{"sector": 20, "multiplier": 1}' http://localhost:8080/games/1/darts
curl -X POST -d '{"sector": 200, "multiplier": 1}' http://localhost:8080/games/1/darts

curl -X GET  http://localhost:8080/games/1
