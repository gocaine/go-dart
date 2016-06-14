package client

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var url = "http://localhost/"
var c = &http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 3,
}

// create the game POST
func startRequest() {
	var body = "serialized JSON"
	resp, err = c.Post(url+games, "application/json", strings.NewReader(body))
	// TODO display the response and save it ?
}
