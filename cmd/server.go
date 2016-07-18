package cmd

import (
	"github.com/gocaine/go-dart/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, arg []string) {
		server := server.NewServer()
		server.Start()
	},
}
