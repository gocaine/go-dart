import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, browserHistory } from 'react-router';
import App from './App';
import ListGames from './ListGames';
import NewGame from './NewGame';
import ViewGame from './ViewGame';
import {addLocaleData, IntlProvider} from 'react-intl';
import en from 'react-intl/locale-data/en';
import fr from 'react-intl/locale-data/fr';
import translates from './locale'

import './index.css';

addLocaleData([...en,...fr]);

const messages = translates[navigator.language.split('-')[0]]


ReactDOM.render((


  <IntlProvider
    locale={navigator.language}
    messages={messages}
    >
    
    <Router history={browserHistory}>
      <Route path="/" component={App} />
      <Route path="newGame" component={NewGame} />
      <Route path="listGames" component={ListGames} />
      <Route path="games/:gameId" component={ViewGame}/>
      <Route path="*" component={App} />
    </Router>
  </IntlProvider>
), document.getElementById('root'));
