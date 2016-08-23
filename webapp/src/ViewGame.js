import React, { Component } from 'react';
import logo from './logo.svg';
import Player from './Player';
import NewPlayerButton from './NewPlayerButton';
import NowPlaying from './NowPlaying';
import PlayerList from './PlayerList';
import NextPlayerButton from './NextPlayerButton'


class ViewGame extends Component {

  constructor(props) {
    super(props)
    this.state = {
      gameId: props.params.gameId,
      Players: []
    }
    console.log("fetching game " + this.state.gameId)
  }

  componentDidMount() {
    console.log("fetching games")
    fetch('/api/games/' + this.state.gameId)
      .then((response) => response.json())
      .then((json) => this.updateGameState(json.game))
      .catch((error) => console.log(error))

    this.ws = new WebSocket('ws://' + window.location.hostname + ':8080/api/games/' + this.state.gameId + '/ws');
    this.ws.onmessage = (event) => this.updateGameState(JSON.parse(event.data))
  }

  updateGameState(game) {
    this.setState({ game: game })
  }

  render() {
    const {gameId} = this.props.params
    const {game} = this.state

    if (!game) {
      return null
    }
    const players = (game.Players ? game.Players.map((player) => <Player key={player.Name} player={player}/>) : [])

    return (
      <div>
        <div className="row">
          <h2><img src={logo} className="App-logo" alt="logo" />Game #{ gameId }</h2>
        </div>
        <div className="row">

          <div className="col s12 l6 push-l6">
            <NowPlaying gameId={ gameId } game={ game } player={ game.Players[game.CurrentPlayer]}/>
          </div>

          <div className="col s12 l5 pull-l6">
            <NewPlayerButton gameId={ gameId }/>
            <PlayerList game={ game } players={ game.Players}/>
          </div>

        </div>
        
      </div>
    );
  }
}

export default ViewGame;