package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// GitHash that was used to build this piece of software
	GitHash = "N/A"
	// BuildDate that was used to build this piece of software
	BuildDate = "N/A"
	// Version is the current project version
	Version = "N/A"
	// ProjectURL is the home page of the software
	ProjectURL = "N/A"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:     v%s\n", Version)
		fmt.Printf("Git hash:    %s\n", GitHash)
		fmt.Printf("Build Date:  %s\n", BuildDate)
		fmt.Println()
		fmt.Printf("Project Url: %s\n", ProjectURL)
	},
}
