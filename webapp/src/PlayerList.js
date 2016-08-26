import React from 'react';

const PlayerList = ({game, players}) => {
  if (game.Ongoing == 4) {
    return (
      <ul className="collection">
        { players.sort((p1, p2) => p1.Rank - p2.Rank).map((p) => <li className="collection-item avatar">
          <i className="material-icons circle">face</i>
          <span className="title"> {p.Name} </span>
          <p>{p.Board}</p>
          <p>{p.Rank}</p>
          <p className="secondary-content">{p.Score}</p>
        </li>
        ) }
      </ul>
      )
  }
  else {
    return (
      <ul className="collection">
        { players.map((p) => <li className="collection-item avatar">
          <i className="material-icons circle">face</i>
          <span className="title"> {p.Name} </span>
          <p>{p.Board}</p>
          <p className="secondary-content">{p.Score}</p>

        </li>
        ) }

      </ul>
    )
  }
}

export default PlayerList;