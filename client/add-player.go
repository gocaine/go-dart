package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func AddPlayer(name string) {
	endpoint := fmt.Sprintf(endpointUrl+"games/%s/user", url.QueryEscape(name))
	resp, err := http.Post(endpoint, "application/json", strings.NewReader("gné"))
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%v\n", resp)
}
