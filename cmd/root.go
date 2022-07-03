package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)
// RootCmd is the root command!
var RootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi, welcome to boilerplate")
	},
}

// Execute runs the main command of the project
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
