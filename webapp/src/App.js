import React, { Component } from 'react';
import { Link } from 'react-router';
import logo from './logo.svg';
import './App.css';
import routes from './config/routes';

class App extends Component {

  constructor(props) {
    super(props)
    console.log("preparing home")
  }

  render() {
    return (
      <div>
        <div className="row">
          <h2><img src={logo} className="App-logo" alt="logo" />Welcome to go-dart</h2>
        </div>
        <div className="row">
          <div className="col s12 m4">
            <div className="icon-block">
              <h2 className="center light-blue-text"><i className="material-icons">add_circle</i></h2>
              <div className="center">
                <Link to={ routes.newGame } className="btn-large waves-effect waves-light orange">Get started</Link>
              </div>
              <p className="light">Create a new game and invite other players to join</p>
            </div>
          </div>
          <div className="col s12 m4">
            <div className="icon-block">
              <h2 className="center light-blue-text"><i className="material-icons">call_merge</i></h2>
              <div className="center">
                <Link to={ routes.listGames } className="btn-large waves-effect waves-light orange">Join existing</Link>
              </div>
              <p className="light">Invite yourself in existing games</p>
            </div>
          </div>
          <div className="col s12 m4">
            <div className="icon-block">
              <h2 className="center light-blue-text"><i className="material-icons">trending_up</i></h2>
              <div className="center">
                <Link to={ routes.home } className="btn-large waves-effect waves-light orange">View statistics</Link>
              </div>
              <p className="light">Explore and analyze statistics of the players</p>
            </div>
          </div>
        </div>
      </div >
    );
  }
}

export default App;
