import React, { Component } from 'react';
import { Link } from 'react-router';
import logo from './logo.svg';

class ListGames extends Component {

  constructor(props) {
    super(props)
    console.log("preparing join")
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
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Join a game</h2>
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