package sip

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/copyleftdev/sipload/internal/rng"
)

// MockRegister simulates sending a SIP REGISTER message before calls.
func MockRegister(ctx context.Context, target string, logger *zap.Logger) error {
	logger.Info("Sending mock REGISTER", zap.String("target", target))

	// pretend it takes 300ms to complete
	select {
	case <-ctx.Done():
		return fmt.Errorf("register canceled")
	case <-time.After(300 * time.Millisecond):
	}

	// 5% chance of failure
	if rng.Intn(20) == 0 {
		logger.Error("Mock register failed", zap.String("target", target))
		return fmt.Errorf("mock register failed")
	}
	return nil
}
