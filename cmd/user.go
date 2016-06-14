package cmd

import (
  "fmt"
  "net/http"
  "net/url"
  "strings"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Add player to existing game",
	Run: func(cmd *cobra.Command, arg []string) {
    name := arg[0]
    endpointUrl := "http://localhost:8080/"
    endpoint := fmt.Sprintf(endpointUrl+"games/%s/user", url.QueryEscape(name))
    resp, err := http.Post(endpoint, "application/json", strings.NewReader("gné"))
    if err != nil {
      fmt.Printf("%s\n", err)
    } else {
      fmt.Printf("%v\n", resp)
    }
	},
}
