// File: internal/sip/simulator_test.go
package sip_test

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"

	"github.com/copyleftdev/sipload/internal/sip"
)

// A helper to get a no-op logger for tests
func getTestLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel) // only log errors or above
	logger, _ := config.Build()
	return logger
}

// TestSimulateCall_Success ensures calls can succeed.
func TestSimulateCall_Success(t *testing.T) {
	logger := getTestLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := sip.SimulateCall(ctx, "sip:mock@fake", "sip:local@fake", logger)
	// The call may fail 10% of the time in the current mock logic
	// so let's do a small loop to reduce flakiness. In real tests, you might
	// stub or mock out RNG or the "fail chance."
	if err != nil {
		t.Logf("First call error: %v. Retrying...", err)
		err = sip.SimulateCall(ctx, "sip:mock@fake", "sip:local@fake", logger)
	}
	assert.Nil(t, err, "SimulateCall should succeed at least once in two tries!")
}

// TestSimulateCall_Canceled verifies we handle context cancellation properly.
func TestSimulateCall_Canceled(t *testing.T) {
	logger := getTestLogger()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately before calling
	cancel()

	err := sip.SimulateCall(ctx, "sip:mock@fake", "sip:local@fake", logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "call canceled")
}
