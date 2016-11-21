# go-dart

[![Build Status](https://travis-ci.org/gocaine/go-dart.svg?branch=master)](https://travis-ci.org/gocaine/go-dart)
[![codecov](https://codecov.io/gh/gocaine/go-dart/branch/master/graph/badge.svg)](https://codecov.io/gh/gocaine/go-dart)


# Contributors

 https://github.com/gocaine/go-dart/graphs/contributors


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

## Accessing the server

The latest stable version of the server is reachable at this address: https://go-dart.mabreizh.fr/


# API

| Method | URI                    | Description                                                        |
|--------|------------------------|--------------------------------------------------------------------|
| GET    | /styles                | Return the list of available game styles                           |
| POST   | /boards                | Register a new board on the server                                 |
| GET    | /boards                | Return the list of boards currently known as 'ALIVE'               |
| GET    | /games                 | Return the list of games                                           |
| POST   | /games                 | Create a new game accordingly to the parameters and returns its id |
| GET    | /games/:gameId         | Return the details of the specified gameId                         |
| DELETE | /games/:gameId         | Delete a game and free the boards                                  |
| POST   | /games/:gameId/players | Add a player to the specified gameId                               |
| POST   | /games/:gameId/hold    | Hold while removing darts from board.                              |
| GET    | /games/:gameId/ws      | Get a WebSocket to listen the events of the specified gameId       |
| POST   | /darts                 | Tell the server a dart hit the board                               |
