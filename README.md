# go-dart

[![Build Status](https://travis-ci.org/gocaine/go-dart.svg?branch=master)](https://travis-ci.org/gocaine/go-dart)
[![codecov](https://codecov.io/gh/gocaine/go-dart/branch/master/graph/badge.svg)](https://codecov.io/gh/gocaine/go-dart)


# Authors

- Guillaume GERBAUD
- Mathieu POUSSE
- Erwann THEBAULT

# Build

First of all, ensure you have checked out the project in `${GOPATH}/github.com/gocaine/go-dart`.

Simply run the command `make binary`

As the build process relies on docker, please first configure it on your local.

If you don't want to use docker, simply run `make dev binary` rely on your local version of `go` and `npm`.

If you want to prepare a binary for ARM, run `make arm binary`

To deploy on the rpi run `make deploy`. In order, you must add the rpi host name and its associated ip to your network configuration (DNS or /etc/hosts)

If you want to build the binary without rebuilding the frontend, run `make binary-noui`. It will ship the previously built version.


## Running on a raspberry pi

After you publish the binary on your rpi (default deploy is in the home directory of the pi user), you can run the following commands: 

 - `./clean-i2c.sh`: this will clean up and reset the GPIO ports direction and status
 - `sudo ./go-dart hardware -b <board-id>`: this will start the software responsible to listenning the hardware and propagating the dart events to the server
 - `./go-dart server`: this starts a game server

# API

- Create a game
  + `POST "/api/games"`
  + return game ID
- Get the current game state
  + `GET "/api/games/{id}"`
  + return a GameState
- Create player
  + `POST "/api/games/{id}/players"`
  + return User ID
- Player state
  + `GET "/api/games/{id}/players/{id}"`
- Dart state
  + `POST "/api/{id}/dart"`
