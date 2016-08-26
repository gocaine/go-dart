import React from 'react';

const Congratulation = ({game, player}) => {
    if (game.Ongoing == 4) {
        return (   
            <div className="card horizontal">
                <div className="card-stacked">
                    <div className="card-content">
                        Congratulations {player.Name}
                    </div>
                </div>
            </div>
            )
    }

}

export default Congratulation;