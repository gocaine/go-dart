import React from 'react';
import NextPlayerButton from './NextPlayerButton'
import ScoreCircle from './ScoreCircle'

const NowPlaying = ({gameId, game, player}) => {
  if (!player) {
    return (<div className="col s12 l6 push-l6">Waiting for players to start</div>)
  }

  return (

    <div className="card horizontal">
      <div className="card-image player-score-card hide-on-small-only">
        <h3>{player.Name}</h3>
        <i className="material-icons circle">face</i>
      </div>
      <div className="card-stacked">
        <div className="card-content">
          <div className="player-score-card-small hide-on-med-and-up">
            <i className="material-icons circle">face</i>
            <span>{player.Name}</span>
          </div>
          <div>
            <ScoreCircle position={0} currentDart={game.CurrentDart} score={player.Visits[0]} />
            <ScoreCircle position={1} currentDart={game.CurrentDart} score={player.Visits[1]} />
            <ScoreCircle position={2} currentDart={game.CurrentDart} score={player.Visits[2]} />
          </div>
        </div>
        <div className="card-action">
          <NextPlayerButton gameId={gameId} game={ game } />
        </div>
      </div>
    </div>

  )
}

export default NowPlaying;