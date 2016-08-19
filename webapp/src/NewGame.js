import React, { Component } from 'react';
import NewGameButton from './NewGameButton'
import logo from './logo.svg';

class NewGame extends Component {

  constructor(props) {
    super(props)
    this.state = {
      flavors: []
    }
  }

  componentDidMount() {
    console.log("fetching flavors")
    fetch('/api/styles')
      .then((response) => response.json())
      .then((json) => this.setState({ flavors: json.styles }))
      .catch((error) => console.log(error))
  }

  render() {
    return (
      <div>
        <div className="row">
          <h2><img src={logo} className="App-logo" alt="logo" />Select flavor</h2>
        </div>
        <div className="row">
          <ul className="collection">
            { this.state.flavors.map((flavor) => <NewGameButton key={flavor.Code} flavor={flavor}/>) }
          </ul>
        </div>
      </div>
    );
  }
}

export default NewGame;