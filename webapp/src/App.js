import React, { Component } from 'react';
import { Link } from 'react-router';
import logo from './logo.svg';
import './App.css';

class App extends Component {

  constructor(props) {
    super(props)
    console.log("preparing home")
  }

  render() {
    return (
      <div>
        <div className="row">

          <div className="center">          
            <img  src={logo} className="App-logo" alt="logo" />
          </div>
          <h4 className="center header-block">Welcome to godart</h4>

          <div className="col s12 m4" >
            <div className="icon-block">
              <div className="center">
                <Link to="/newGame" className="btn-large waves-effect waves-light btn light-blue"><i className="material-icons left large ">add</i>New Game</Link>
                <p className="light hide-on-small-only">Create a new game and invite other players to join</p>
              </div>
            </div>            
          </div>

          <div className="col s12 m4">
            <div className="icon-block">
              <div className="center">
                <Link to="/newGame" className="btn-large waves-effect waves-light btn light-blue"><i className="material-icons left large">call_merge</i>Join existing</Link>
                <p className="light hide-on-small-only">Invite yourself in existing games</p>
              </div>
            </div>
          </div>

          <div className="col s12 m4">
            <div className="icon-block">
              <div className="center">
                <Link to="/newGame" className="btn-large waves-effect waves-light btn light-blue"><i className="material-icons left large">subscriptions</i>View statistics</Link>
                <p className="hide-on-small-only light">Explore and analyze statistics of the players</p>
              </div>
            </div>
          </div>
        </div>
      </div >
    );
  }
}

export default App;
