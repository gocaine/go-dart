import React from 'react';

const ScoreCircle = ({ position, currentDart, score  }) => {

    if (score != undefined && (score.Pos == 2 || score.Pos == 3)) {

        return (
            <div className="score-circle-wrapper">
                <div className="circle score-circle valign-wrapper"
                    style={{ 'background-color': currentDart == position ? '#ef6c00' : (currentDart > position ? 'orange' : 'grey') } }>
                    <div className="valign center-align score-circle-content">
                        {  score.Val  }
                    </div>
                </div>
                <div className="circle bubble">x{score.Pos}</div>
            </div>
        )
    }
    else {
        return (
            <div className="circle score-circle valign-wrapper"
                style={{ 'background-color': currentDart == position ? '#ef6c00' : (currentDart > position ? 'orange' : 'grey') } }>
                <div className="valign center-align score-circle-content">
                    {  score != undefined ? score.Val : ''  }
                </div>
            </div>

        )
    }



}

export default ScoreCircle;