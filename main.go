package main

import (
	"os"

	"github.com/copyleftdev/sipload/cmd"
)

// main simply executes the root command.
func main() {
	rootCmd := cmd.GetRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
