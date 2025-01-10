package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show sipload version",
	Run: func(cmd *cobra.Command, args []string) {
		// You might dynamically set this at build time with ldflags.
		fmt.Println("sipload version v0.2.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
