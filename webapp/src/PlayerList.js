import React from 'react';

const PlayerList = ({game, players}) => {
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

export default PlayerList;