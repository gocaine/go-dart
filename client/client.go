package client

import (
	"bytes"
	"fmt"
	"go-dart/common"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/spf13/viper"
)

type DartClient struct {
	base *sling.Sling
}

func NewClient() *DartClient {

	var c = &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 3,
	}

	var endpointURL = viper.GetString("server")
	client := DartClient{}
	client.base = sling.New().Base(endpointURL).Client(c)
	return &client
}

func (client *DartClient) CreateGame(style string) (response GameResponse, err error) {
	var failure Failure
	game := common.GameRepresentation{Style: style}
	rawResponse, err := client.base.New().Post("games").BodyJSON(game).Receive(&response, &failure)
	err = client.manageFailure(rawResponse, failure, err)
	return
}

func (client *DartClient) GetState(gameId int) (response GameResponse, err error) {
	var failure Failure
	path := fmt.Sprintf("games/%d", gameId)
	rawResponse, err := client.base.New().Get(path).Receive(&response, &failure)
	err = client.manageFailure(rawResponse, failure, err)
	return
}

func (client *DartClient) CreatePlayer(gameId int, name string) (resource string, err error) {
	var failure Failure
	player := common.PlayerRepresentation{Name: name}

	path := fmt.Sprintf("games/%d/players", gameId)
	rawResponse, err := client.base.New().Post(path).BodyJSON(player).Receive(&resource, &failure)
	err = client.manageFailure(rawResponse, failure, err)
	return
}

func (client *DartClient) FireDart(gameId int, sector int, multiplier int) (state common.GameState, err error) {
	var failure Failure
	dartRep := common.DartRepresentation{Sector: sector, Multiplier: multiplier}

	path := fmt.Sprintf("games/%d/darts", gameId)
	rawResponse, err := client.base.New().Post(path).BodyJSON(dartRep).Receive(&state, &failure)
	err = client.manageFailure(rawResponse, failure, err)
	return
}

//manageFailure convert server failure or error into error
// if no failure is returned and no 200 HTTP code, an error is returned
//FIXME manage other good HTTP status code
func (client *DartClient) manageFailure(rawResponse *http.Response, failure Failure, err error) error {
	if rawResponse.StatusCode >= 400 && failure.Status == "" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(rawResponse.Body)
		return fmt.Errorf("Technical error -> http code  %d, response %s ", rawResponse.StatusCode, buf.String())
	} else if failure.Status != "" {
		return fmt.Errorf("Status : %s, Error :  %s", failure.Status, failure.Error)
	}
	return nil
}
