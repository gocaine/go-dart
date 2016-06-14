package cmd

import (
	"fmt"
        "go-dart/client"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "go-dart",
	Short: "go-dart is cool",
	Long: `A better dart game than the chinese one.
	Complete doc at voir http://github.com/Zenika/go-dart.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(addPlayerCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Begin a new game",
	Long:  `Create a new game`,
	Run: func(cmd *cobra.Command, arg []string) {
		// API CALL
	},
}

var addPlayerCmd = &cobra.Command{
	Use:   "add-player",
	Short: "Add player to existing game",
	Run: func(cmd *cobra.Command, arg []string) {
            client.AddPlayer(arg[0])
	},
}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  `Display current version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-dart v0.0 HACKATON")
	},
}
