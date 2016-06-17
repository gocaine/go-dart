package cmd

import (
	"fmt"
	"go-dart/client"
	"strconv"

	"github.com/spf13/cobra"
)

var addPlayerCmd = &cobra.Command{
	Use:   "addplayer <gameID>",
	Short: "Add player to existing game",
	Run: func(cmd *cobra.Command, arg []string) {
		if len(arg) < 1 {
			fmt.Println("Missing argument. Please provide game ID!")
		} else {
			gameID := arg[0]
			if i, err := strconv.Atoi(gameID); err == nil {
				AddPlayer(i)
			} else {
				fmt.Println("Provided <gameID> must be an integer")
			}
		}
	},
}

// AddPlayer to the game
func AddPlayer(gameID int) {
	resp, err := client.Get("game/" + strconv.Itoa(gameID) + "/user")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("%s\n", string(resp))
}
