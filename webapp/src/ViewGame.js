import React, { Component } from 'react';
import logo from './logo.svg';
import Player from './Player';
import NewPlayerButton from './NewPlayerButton';
import NowPlaying from './NowPlaying';
import PlayerList from './PlayerList';
import NextPlayerButton from './NextPlayerButton'
import Congratulation from './Congratulation'


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

    var wsProtocol = "wss";
    if (window.location.protocol == "http:") {
      // keep using non secured connection
      wsProtocol = "ws";
    }
    this.ws = new WebSocket(wsProtocol + '://' + window.location.host + '/api/games/' + this.state.gameId + '/ws');
    this.ws.onmessage = (event) => this.handleWsMessage(JSON.parse(event.data))
  }

  handleWsMessage(message) {
    if(message.Kind === 'status') {
      this.setState({game: message.Data})
    }
  }

  render() {
    const {gameId} = this.props.params
    const {game} = this.state

    if (!game) {
      return null
    }
    const players = (game.Players ? game.Players.map((player) => <Player key={player.Name} player={player}/>) : [])

    // Display NewPlayerButton only if ongoing is INITIALIZING or READY
    const newPlayerButton = (game.Ongoing <= 1 ? <NewPlayerButton gameId={ gameId }/> : "")


    const playerPanel = (game.Ongoing == 4 ?
      <Congratulation game={ game } player={ game.Players[game.CurrentPlayer]} > </Congratulation> :
      <NowPlaying gameId={ gameId } game={ game } player={ game.Players[game.CurrentPlayer]}/>)

    return (
      <div>
        <div className="row">
          <h2 className="col s12 l6 offset-l3"><img src={logo} className="App-logo" style={{'vertical-align': 'middle'}} alt="logo" /><span style={{'vertical-align': 'middle'}} >Game #{ gameId }</span></h2>
        </div>
        <div className="row">

          <div className="col s12 l6 offset-l3">
            {newPlayerButton}
          </div>

          <div className="col s12 l6 offset-l3">
            {playerPanel}
          </div>

          <div className="col s12 l6 offset-l3">
            <PlayerList game={ game } players={ game.Players}/>
          </div>

        </div>

      </div>
    );
  }
}

export default ViewGame;
