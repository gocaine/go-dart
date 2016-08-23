import React, { Component } from 'react';
import { browserHistory } from 'react-router';
import {Input} from 'react-materialize'



class NewPlayerButton extends Component {

  constructor(props) {
    super(props)
    this.state = {
      boards: []
    }
  }

  componentDidMount() {
    //TODO use a react library for materialyzecss
    this.listBoards()
  }

  addPlayer() {
    fetch('/api/games/' + this.props.gameId + "/players", {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: this._nameinput.value,
        board: this._boardinput.selectInput.value
      })
    })
      .then((response) => response.json())
      .then((json) => console.log(json))
      .catch((error) => console.log(error))
  }

  listBoards() {
    fetch('/api/boards')
      .then(response => response.json())
      .then((json) => this.setState({ boards: json }))
      .catch((error) => console.log(error))
  }



  render() {
    return (
      <div className="col s12">
        <div className="input-field col s12 l5">
          <input  ref={(c) => this._nameinput = c} id="playerName" type="text" className="validate"/>
          <label htmlFor="playerName">Player name</label>
        </div>
        <div className=" col s8 l4">
          <Input type='select' label="Board" ref={(c) => this._boardinput = c}>
            { this.state.boards.map((board) => <option value={board}>{ board }</option>) }
          </Input>
        </div>
        <div className="input-field col s4 l3">
          <a className="waves-effect waves-light btn light-blue" onClick={() => this.addPlayer() }>Add</a>
        </div>
      </div>
      )
  }
}

export default NewPlayerButton;