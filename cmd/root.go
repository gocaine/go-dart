package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  `Display current version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-dart v0.0 HACKATON")
	},
}
