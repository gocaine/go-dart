package client

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var endpointUrl = "http://localhost:8080/"
var c = &http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 3,
}

// create the game POST
func startRequest() {
	var body = "serialized JSON"
	resp, err := c.Post(endpointUrl+"games", "application/json", strings.NewReader(body))
	// TODO display the response and save it ?
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%v\n", resp)
}
