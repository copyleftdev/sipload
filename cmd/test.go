package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/copyleftdev/sipload/internal/load"
	"github.com/copyleftdev/sipload/internal/stats"
)

// Flags (overrides)
var (
	flagTargetURI      string
	flagCallsPerSecond float64
	flagConcurrency    int
	flagTestDuration   time.Duration
	flagLocalContact   string
	flagRegisterFirst  bool
)

// testCmd is the subcommand for running the SIP load test.
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run a SIP load test with specified parameters",
	RunE:  runTest,
}

func init() {
	// Define flags for the test command
	testCmd.Flags().StringVar(&flagTargetURI, "target", "", "SIP target URI (overrides config)")
	testCmd.Flags().Float64Var(&flagCallsPerSecond, "calls-per-second", 0, "Call generation rate (CPS)")
	testCmd.Flags().IntVar(&flagConcurrency, "concurrency", 0, "Max number of simultaneous calls")
	testCmd.Flags().DurationVar(&flagTestDuration, "duration", 0, "How long to run the test (0 = infinite/until Ctrl+C)")
	testCmd.Flags().StringVar(&flagLocalContact, "contact", "", "Local Contact URI")
	testCmd.Flags().BoolVar(&flagRegisterFirst, "register-first", false, "Send REGISTER before placing calls")

	// Add the test command to the root
	rootCmd.AddCommand(testCmd)
}

// runTest is the function that actually runs when "sipload test" is called.
func runTest(cmd *cobra.Command, args []string) error {
	// Gather configuration using Viper + flags override
	cfg := getTestConfig()

	logger.Info("[INFO] Starting SIP load test",
		zap.String("targetURI", cfg.TargetURI),
		zap.Float64("callsPerSecond", cfg.CallsPerSecond),
		zap.Int("concurrency", cfg.Concurrency),
		zap.Duration("duration", cfg.TestDuration),
		zap.String("localContact", cfg.LocalContact),
		zap.Bool("registerFirst", cfg.RegisterFirst),
	)
	fmt.Println("Press Ctrl+C to stop early...")

	// Handle graceful shutdown on Ctrl+C or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	// Prepare stats collector
	statsCollector := stats.NewCollector()

	// Prepare rate limiter
	burst := cfg.Concurrency * 2
	if burst < 1 {
		burst = 1
	}
	limiter := rate.NewLimiter(rate.Limit(cfg.CallsPerSecond), burst)

	// Create and run the tester
	tester := load.NewTester(cfg, limiter, statsCollector, logger)
	if err := tester.Run(ctx); err != nil {
		logger.Error("[ERROR] Tester encountered an error", zap.Error(err))
	}

	// Print final stats
	printResults(statsCollector)
	return nil
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Info("[INFO] Received stop signal... Shutting down.")
	cancel()
}

func printResults(collector *stats.Collector) {
	totalCalls := collector.TotalCalls()
	totalFailures := collector.TotalFailures()
	elapsed := collector.Elapsed()

	fmt.Println("\n--- Test Results ---")
	fmt.Printf("Total Calls Attempted: %d\n", totalCalls)
	fmt.Printf("Total Failures:        %d\n", totalFailures)
	fmt.Printf("Test Duration:         %s\n", elapsed)

	if totalCalls > 0 {
		successCalls := totalCalls - totalFailures
		successRate := float64(successCalls) / float64(totalCalls) * 100
		fmt.Printf("Success Rate:          %.2f%%\n", successRate)
	}
	fmt.Println("[INFO] Load test completed.")
}

// getTestConfig merges config from Viper with any CLI flags.
func getTestConfig() *load.TestConfig {
	// The approach: first read from config file, then override with CLI flags if set.

	cfg := &load.TestConfig{
		TargetURI:      viper.GetString("target_uri"),
		CallsPerSecond: viper.GetFloat64("calls_per_second"),
		Concurrency:    viper.GetInt("concurrency"),
		TestDuration:   time.Second * time.Duration(viper.GetInt("duration")),
		LocalContact:   viper.GetString("local_contact"),
		RegisterFirst:  viper.GetBool("register_first"),
	}

	// CLI overrides
	if flagTargetURI != "" {
		cfg.TargetURI = flagTargetURI
	}
	if flagCallsPerSecond > 0 {
		cfg.CallsPerSecond = flagCallsPerSecond
	}
	if flagConcurrency > 0 {
		cfg.Concurrency = flagConcurrency
	}
	if flagTestDuration > 0 {
		cfg.TestDuration = flagTestDuration
	}
	if flagLocalContact != "" {
		cfg.LocalContact = flagLocalContact
	}
	if flagRegisterFirst {
		cfg.RegisterFirst = true
	}

	return cfg
}
