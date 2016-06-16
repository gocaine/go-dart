package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Begin a new game",
	Run: func(cmd *cobra.Command, arg []string) {
		// API CALL
	},
}
