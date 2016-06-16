package cmd

import (
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

// Start a new game
func Start(style string) {
	client.Post("games", "{\"style\": \""+style+"\"}")
}
