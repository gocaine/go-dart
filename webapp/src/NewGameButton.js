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
          <span>{this.props.flavor.Code}</span>
          <a onClick={() => this.newGame(this.props.flavor.Code) } className="secondary-content"><i className="material-icons light-blue-text">keyboard_arrow_right</i></a>
        </div>
        <div className="collapsible-body"><p>Voluptate ea aliquip esse consequat eu reprehenderit laborum sunt sit. Esse labore duis amet sint in veniam aute esse enim. Adipisicing culpa quis aliqua est excepteur magna. Nostrud amet incididunt irure duis ea exercitation qui. Est pariatur est pariatur non pariatur anim dolore velit reprehenderit commodo id consequat proident.</p></div>
      </li>
    )
  }
}

export default NewGameButton;