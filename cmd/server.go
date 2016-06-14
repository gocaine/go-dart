package cmd

import (
	"github.com/spf13/cobra"
  "go-dart/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, arg []string) {
		server := server.NewServer()
		server.Start()
	},
}
