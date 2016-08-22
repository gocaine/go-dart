import React, {Component} from 'react'

class NextPlayerButton extends Component {

    nextPlayer() {
        fetch('/api/games/' + this.props.gameId + '/skip', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({})
        })
            .then((response) => response.json())
            .catch((error) => console.log(error))
    }

    render() {
        return (
            <div>
                <a className="waves-effect waves-light btn light-blue" onClick={() => this.nextPlayer() }>Next</a>
            </div>
        )
    }
}

export default NextPlayerButton