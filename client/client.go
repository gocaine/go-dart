package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gocaine/go-dart/common"

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
	client.base = sling.New().Base(endpointURL + "/api/").Client(c)
	return &client
}

func (client *DartClient) FireDart(board string, sector int, multiplier int) (state common.GameState, err error) {
	var failure Failure
	dartRep := common.DartRepresentation{Board: board, Sector: sector, Multiplier: multiplier}

	rawResponse, err := client.base.New().Post("darts").BodyJSON(dartRep).Receive(&state, &failure)
	err = client.manageFailure(rawResponse, failure, err)
	return
}

//manageFailure convert server failure or error into error
// if no failure is returned and no 200 HTTP code, an error is returned
//FIXME manage other good HTTP status code
func (client *DartClient) manageFailure(rawResponse *http.Response, failure Failure, err error) error {
	if rawResponse != nil && rawResponse.StatusCode >= 400 && failure.Status == "" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(rawResponse.Body)
		return fmt.Errorf("Technical error -> http code  %d, response %s ", rawResponse.StatusCode, buf.String())
	} else if failure.Status != "" {
		return fmt.Errorf("Status : %s, Error :  %s", failure.Status, failure.Error)
	} else {
		return err
	}
}
