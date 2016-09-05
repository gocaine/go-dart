import React, { Component } from 'react';
import { browserHistory } from 'react-router';
import {Input} from 'react-materialize'



class NewPlayerButton extends Component {

  constructor(props) {
    super(props)
    this.state = {
      boards: []
    }

    this.handleBoardChange = this.handleBoardChange.bind(this)
    this.handlePlayerNameChange = this.handlePlayerNameChange.bind(this)
  }

  componentDidMount() {
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
        name: this.name,
        board: this.board
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


  handlePlayerNameChange(event) {
    this.name = event.target.value
  }

  handleBoardChange(event) {
    this.board = event.target.value
  }


  render() {
    return (
      <div className="col s12">
        <div className="input-field col s12 l5">
          <input id="playerName" type="text" className="validate" onChange={this.handlePlayerNameChange}/>
          <label htmlFor="playerName">Player name</label>
        </div>
        <div className=" col s8 l4">
          <Input type='select' label="Board"  onChange={this.handleBoardChange}>
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