import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, browserHistory } from 'react-router';
import App from './App';
import ListGames from './ListGames';
import NewGame from './NewGame';
import ViewGame from './ViewGame';
import './index.css';
import { createHistory, useBasename } from 'history'

// Run our app under the /web URL (even in dev mode).
const history = useBasename(createHistory)({
  basename: '/web'
})

ReactDOM.render((
  <Router history={ history }>
    <Route path="/" component={App}/>
    <Route path="newGame" component={NewGame}/>
    <Route path="listGames" component={ListGames}/>
    <Route path="games/:gameId" component={ViewGame}/>
    <Route path="*" component={App}/>
  </Router>
), document.getElementById('root'));
