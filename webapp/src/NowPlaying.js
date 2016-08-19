import React from 'react';

const NowPlaying = ({game, player}) => {
  if (!player) {
    return (<div>Waiting for players to start</div>)
  }
  return (
    <div>
      <div>Now playing: { player.Name } - Score is { player.Score }</div>
      <div>Throwing #{ game.CurrentDart }...</div>
      <div>Last dart hit { game.LastSector.Pos } * {game.LastSector.Val}</div>
    </div>
  )
}

export default NowPlaying;