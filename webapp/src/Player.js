import React, { Component } from 'react';

const Player = ({player}) => {
  return (
    <div key={player.Name}>Player: {player.Name} - Score is {player.Score}</div>
  )
}

export default Player;