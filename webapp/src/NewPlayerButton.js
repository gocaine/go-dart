import React, { Component } from 'react';

class NewPlayerButton extends Component {

  addPlayer() {
    fetch('/api/games/' + this.props.gameId + "/players", {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: this.refs.name.value,
        board: "Rennes"
      })
    })
      .then((response) => response.json())
      .then((json) => console.log(json))
      .catch((error) => console.log(error))
  }

  render() {
    return (
      <div className="row">
        <div className="input-field col s6">
          <input  ref="name" id="playerName" type="text" className="validate"/>
          <label htmlFor="playerName">Player name</label>
        </div>
        <div className="input-field col s6">
          <a className="waves-effect waves-light btn light-blue" onClick={() => this.addPlayer() }>Add</a>
        </div>
      </div>)
  }
}

export default NewPlayerButton;