package client

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var c = &http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       time.Second * 3,
}

// Post request on API endpoint and return the answer
func Post(cmd, body string) (*http.Response, error) {
	var endpointURL = "http://" + viper.GetString("server") + "/"
	resp, err := c.Post(endpointURL+cmd, "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	fmt.Printf("%v\n", resp)
	return resp, nil
}

// Get request on API endpoint and return the answer
func Get(cmd, body string) (*http.Response, error) {
	var endpointURL = "http://" + viper.GetString("server") + "/"
	resp, err := c.Get(endpointURL + cmd)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	fmt.Printf("%v\n", resp)
	return resp, nil
}
