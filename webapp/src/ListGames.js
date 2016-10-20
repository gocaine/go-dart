import React, { Component } from 'react';
import { Link } from 'react-router';
import logo from './logo.svg';
import {FormattedMessage} from 'react-intl';

class ListGames extends Component {

  constructor(props) {
    super(props)
    this.state = {
      games: []
    }
  }

  componentDidMount() {
    console.log("fetching games")
    fetch('/api/games')
      .then((response) => response.json())
      .then((json) => this.setState({games: json}))
      .catch((error) => console.log(error))
  }

  render() { 
   return (
      <div className="App">
        <div className="App-header">
          <img style={{'vertical-align': 'middle'}} src={logo} className="App-logo" alt="logo" />
          <h2 style={{'vertical-align': 'middle'}} ><FormattedMessage id='listGame.join' defaultMessage='Join a game'/></h2>
        </div>
        <ul>
        {
          this.state.games.map((game) => <li key={game}><Link to={`/games/${game}`}>Game #{game}</Link></li>)
        }
        </ul>
      </div>
    );
  }
}

export default ListGames;