package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"

	log "github.com/Sirupsen/logrus"
)

// RootCmd is the default command
var RootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", cmd.UsageString())
		os.Exit(-1)
	},
}

func init() {
	RootCmd.AddCommand(serverCmd())
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(hardwareCmd())

	RootCmd.PersistentFlags().StringP("log", "l", "info", "Log level (debug|info|warn|error)")
	viper.BindPFlag("log", RootCmd.PersistentFlags().Lookup("log"))
	cobra.OnInitialize(configureLogger)
}

func configureLogger() {
	level, err := log.ParseLevel(viper.GetString("log"))
	if err != nil {
		log.Panicf("invalid log level %s", level)
	}
	log.SetLevel(level)
	log.Debugf("Log level is now set at %s", level)
}
