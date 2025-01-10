package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// rootCmd is the base command of the CLI.
var rootCmd = &cobra.Command{
	Use:   "sipload",
	Short: "A token-bucket-based SIP load-testing tool in Go",
	Long: `sipload is a CLI tool that uses a token-bucket algorithm to
generate SIP-like traffic at a controlled rate with concurrency limits.`,
}

// logger is a global zap.Logger to be shared across commands.
// A real project might pass this via dependency injection.
var logger *zap.Logger

// GetRootCmd returns the root command for main.go to execute.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// 1. Initialize Viper
	initConfig()

	// 2. Initialize logger (zap)
	initLogger()

	// You can add any persistent flags here if needed, e.g.:
	// rootCmd.PersistentFlags().String("log-level", "info", "Log level")

	// Example of reading that flag into Viper
	// viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
}

func initConfig() {
	// Tell Viper where to look for the config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs") // folder to look for config.yaml
	viper.AddConfigPath(".")         // also look in current directory

	// Read env variables with prefix "SIPLOAD_"
	viper.SetEnvPrefix("SIPLOAD")
	viper.AutomaticEnv()

	// Try to read the config file; ignore error if not found
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("[WARN] Could not read config file: %v\n", err)
	}
}

func initLogger() {
	// Create a default production logger; for local dev, you might prefer zap.NewDevelopment()
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
	}
}
