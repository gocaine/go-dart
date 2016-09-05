package cmd

import (
	"github.com/gocaine/go-dart/server"

	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Run: func(cmd *cobra.Command, arg []string) {

			flag := cmd.Flag("port")
			port := flag.Value.String()

			server := server.NewServer()
			server.Start(port)
		},
	}

	serverCmd.Flags().IntP("port", "p", 8080, "listening port")

	return serverCmd
}
