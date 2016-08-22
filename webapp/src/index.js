import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, browserHistory } from 'react-router';
import App from './App';
import ListGames from './ListGames';
import NewGame from './NewGame';
import ViewGame from './ViewGame';
import routes from './config/routes';
import './index.css';

ReactDOM.render((
  <Router history={ browserHistory }>
    <Route path={ routes.home } component={ App }/>
    <Route path={ routes.newGame } component={ NewGame }/>
    <Route path={ routes.listGames } component={ ListGames }/>
    <Route path={ routes.joinGameById } component={ ViewGame }/>
    <Route path="*" component={ App }/>
  </Router>
), document.getElementById('root'));
