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

// POST request on API endpoint and return the answer
func Request(cmd, body string) (*http.Response, error) {
	resp, err := c.Post(endpointUrl+cmd, "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	fmt.Printf("%v\n", resp)
	return resp, nil
}
