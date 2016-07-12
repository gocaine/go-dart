package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// RootCmd is the default command
var RootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(-1)
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.AddCommand(sampleCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(hardwareCmd())
	RootCmd.PersistentFlags().StringP("server", "s", "http://localhost:8080/", "Server address")
	viper.BindPFlag("server", RootCmd.PersistentFlags().Lookup("server"))
}
