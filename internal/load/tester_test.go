// File: internal/load/tester_test.go
package load_test

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"

	"github.com/stretchr/testify/assert"

	"github.com/copyleftdev/sipload/internal/load"
	"github.com/copyleftdev/sipload/internal/stats"
)

// A quiet logger for testing
func getQuietLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	logger, _ := config.Build()
	return logger
}

func TestTester_Run(t *testing.T) {
	logger := getQuietLogger()
	cfg := &load.TestConfig{
		TargetURI:      "sip:mock@fake",
		CallsPerSecond: 10,
		Concurrency:    2,
		TestDuration:   1 * time.Second,
		LocalContact:   "sip:test@local",
		RegisterFirst:  true, // might do a REGISTER call
	}

	limiter := rate.NewLimiter(rate.Limit(cfg.CallsPerSecond), cfg.Concurrency*2)
	collector := stats.NewCollector()
	tester := load.NewTester(cfg, limiter, collector, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := tester.Run(ctx)
	assert.NoError(t, err, "Tester.Run should complete without error")

	// Verify some calls happened
	totalCalls := collector.TotalCalls()
	t.Logf("Total calls: %d", totalCalls)
	assert.Greater(t, totalCalls, 0, "Expected at least 1 call")
}
