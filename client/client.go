package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gocaine/go-dart/common"

	log "github.com/Sirupsen/logrus"
	"github.com/dghubble/sling"
)

// DartClient the API client
type DartClient struct {
	base       *sling.Sling
	board      string
	tickerQuit chan bool
}

// NewClient the client constructor
func NewClient(endpointURL string, board string) *DartClient {

	var c = &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 3,
	}

	client := DartClient{board: board}
	client.base = sling.New().Base(endpointURL + "/api/").Client(c)
	client.healthCheckInit()
	return &client
}

func (client *DartClient) healthCheckInit() {
	client.ping()
	ticker := time.NewTicker(common.HealthCheckDelay)
	client.tickerQuit = make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				client.ping()
			case <-client.tickerQuit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Shutdown releases the resources held by the client
func (client *DartClient) Shutdown() {
	client.tickerQuit <- true
}

func (client *DartClient) ping() {
	log.Debug("ping...")
	_, err := client.base.New().Post("boards").BodyJSON(common.BoardRepresentation{Name: client.board}).ReceiveSuccess(nil)
	if err != nil {
		// server is down ?
		log.Panicf("error during ping: %v", err)
	}
}

// FireDart sends a Dart to the server
func (client *DartClient) FireDart(sector int, multiplier int) (state common.GameState, err error) {
	var failure Failure
	dartRep := common.DartRepresentation{Board: client.board, Sector: sector, Multiplier: multiplier}

	rawResponse, err := client.base.New().Post("darts").BodyJSON(dartRep).Receive(&state, &failure)
	err = client.handleFailure(rawResponse, failure, err)
	return
}

//handleFailure converts server failure or error into error
// if no failure is returned and no 200 HTTP code, an error is returned
//FIXME manage other good HTTP status code
func (client *DartClient) handleFailure(rawResponse *http.Response, failure Failure, err error) error {
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
