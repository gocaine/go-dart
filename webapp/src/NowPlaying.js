import React from 'react';
import NextPlayerButton from './NextPlayerButton'


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
              <div className="circle score-circle valign-wrapper" 
                  style={{'background-color': game.CurrentDart==0 ? '#ef6c00':'orange'} }>
                <div className="valign center-align score-circle-content">
                  { game.CurrentDart==1 ? game.LastSector.Pos * game.LastSector.Val : '' }
                </div>
              </div>
              <div className="circle score-circle valign-wrapper" 
                  style={{'background-color': game.CurrentDart==1 ? '#ef6c00' : (game.CurrentDart>1 ? 'orange':'grey') } }>
                <div className="valign center-align score-circle-content">
                  { game.CurrentDart==2 ? game.LastSector.Pos * game.LastSector.Val : '' }
                </div>
              </div>
              <div className="circle score-circle valign-wrapper" 
                  style={{'background-color': game.CurrentDart==2 ? '#ef6c00' : 'grey'} }>
                <div className="valign center-align score-circle-content">
                  { game.CurrentDart==3 ? game.LastSector.Pos * game.LastSector.Val : '' }
                </div>
              </div>
            </div>
          </div>
          <div className="card-action">
            <NextPlayerButton gameId={ gameId } />
          </div>
        </div>
      </div>
     
  )
}

export default NowPlaying;