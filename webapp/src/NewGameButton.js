import React, { Component } from 'react';
import { browserHistory } from 'react-router';
import routes from './config/routes';

class NewGameButton extends Component {

  newGame(flavor) {
    fetch('/api/games', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        style: flavor
      })
    })
    .then((response) => response.json())
    .then((json) => browserHistory.push(routes.joinGameById.replace(':gameId', json.id)))
    .catch((error) => console.log(error))    
  }

  render() {
   return (
     <li className="collection-item">
     <span>{ this.props.flavor.Code }</span>
     <a onClick={ () => this.newGame(this.props.flavor.Code) } className="secondary-content"><i className="material-icons light-blue-text">keyboard_arrow_right</i></a>
     </li>)
  }
}

export default NewGameButton;