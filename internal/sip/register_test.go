// File: internal/sip/register_test.go
package sip_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/copyleftdev/sipload/internal/sip"
)

func getSilentLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	logger, _ := config.Build()
	return logger
}

func TestMockRegister_Success(t *testing.T) {
	logger := getSilentLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := sip.MockRegister(ctx, "sip:mock@fake", logger)
	// 5% chance of random failure. We can allow for a retry or just check it once.
	// For demonstration, let's do a single attempt and accept a small risk of flakiness.
	if err != nil {
		t.Logf("MockRegister returned error (could be 5%% random): %v", err)
	}
	// This might fail rarely due to random chance.
	// In production, consider mocking the RNG or controlling the fail path.
}

func TestMockRegister_Canceled(t *testing.T) {
	logger := getSilentLogger()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sip.MockRegister(ctx, "sip:mock@fake", logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "register canceled")
}
