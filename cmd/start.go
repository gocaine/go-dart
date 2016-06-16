package cmd

import (
	"encoding/json"
	"fmt"
	"go-dart/client"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <gamestyle>",
	Short: "Begin a new game, see doc for gamestyle",
	Run: func(cmd *cobra.Command, arg []string) {
		if len(arg) < 1 {
			fmt.Println("Missing argument. Retry with a gamestyle : 101, 301...")
		} else {
			Start(arg[0])
		}
	},
}

type game struct {
	ID int
}

// Start a new game and display the corresponding ID
func Start(style string) {
	resp, err := client.Post("games", "{\"style\": \""+style+"\"}")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	var respJSON game
	err = json.Unmarshal(resp, &respJSON)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Game id: %d\n", respJSON.ID)
}
