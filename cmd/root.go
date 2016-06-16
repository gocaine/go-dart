package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the default command
var RootCmd = &cobra.Command{
	Use:   "go-dart",
	Short: "go-dart is cool",
	Long: `A better dart game than the chinese one.
	Complete doc at voir http://github.com/Zenika/go-dart.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(userCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.PersistentFlags().StringP("server", "s", "localhost:8080", "Server address")
	viper.BindPFlag("server", RootCmd.PersistentFlags().Lookup("server"))
}
