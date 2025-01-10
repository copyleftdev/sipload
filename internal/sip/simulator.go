package sip

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/copyleftdev/sipload/internal/rng"
)

// SimulateCall is a placeholder for actual SIP call logic.
// Replace with a real SIP library in production.
func SimulateCall(ctx context.Context, targetURI, contact string, logger *zap.Logger) error {
	// Simulate call setup time ~100-400ms
	setupTime := time.Millisecond * time.Duration(100+rng.Intn(300))

	select {
	case <-ctx.Done():
		logger.Warn("Call canceled by context")
		return fmt.Errorf("call canceled")
	case <-time.After(setupTime):
	}

	// 10% chance of failure
	if rng.Intn(10) == 0 {
		logger.Error("Mock call failed", zap.String("targetURI", targetURI))
		return fmt.Errorf("mock call failed")
	}

	logger.Debug("Call succeeded", zap.String("targetURI", targetURI))
	return nil
}
