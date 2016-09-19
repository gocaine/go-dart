import React, { Component } from 'react';
import { browserHistory } from 'react-router';

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
      .then((json) => browserHistory.push("/games/" + json.id))
      .catch((error) => console.log(error))
  }

  render() {
    return (
      <li>
        <div className="collapsible-header">
          <span>{this.props.flavor.Name}</span>
          <a onClick={() => this.newGame(this.props.flavor.Code) } className="secondary-content"><i className="material-icons light-blue-text">keyboard_arrow_right</i></a>
        </div>
        <div className="collapsible-body"><p>{this.props.flavor.Rules}</p></div>
      </li>
    )
  }
}

export default NewGameButton;
