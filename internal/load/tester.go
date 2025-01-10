package load

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/copyleftdev/sipload/internal/sip"
	"github.com/copyleftdev/sipload/internal/stats"
)

// TestConfig holds the parameters for a SIP load test.
type TestConfig struct {
	TargetURI      string
	CallsPerSecond float64
	Concurrency    int
	TestDuration   time.Duration
	LocalContact   string
	RegisterFirst  bool
}

// Tester orchestrates the SIP load test.
type Tester struct {
	cfg       *TestConfig
	limiter   *rate.Limiter
	collector *stats.Collector
	logger    *zap.Logger
}

// NewTester creates a new Tester.
func NewTester(cfg *TestConfig, limiter *rate.Limiter, collector *stats.Collector, logger *zap.Logger) *Tester {
	return &Tester{
		cfg:       cfg,
		limiter:   limiter,
		collector: collector,
		logger:    logger,
	}
}

// Run executes the SIP load test using concurrency, rate-limiting, and stats.
func (t *Tester) Run(ctx context.Context) error {
	// Possibly do a REGISTER first
	if t.cfg.RegisterFirst {
		err := sip.MockRegister(ctx, t.cfg.TargetURI, t.logger)
		if err != nil {
			t.logger.Error("[ERROR] REGISTER step failed", zap.Error(err))
		} else {
			t.logger.Info("[INFO] REGISTER completed successfully")
		}
	}

	startTime := time.Now()
	t.collector.StartTimer()

	// If TestDuration > 0, we set an end time; else run until canceled.
	var endTime time.Time
	if t.cfg.TestDuration > 0 {
		endTime = startTime.Add(t.cfg.TestDuration)
	}

	// Concurrency control via channel semaphore
	sem := make(chan struct{}, t.cfg.Concurrency)
	var wg sync.WaitGroup

LOOP:
	for {
		select {
		case <-ctx.Done():
			// Test was cancelled
			break LOOP
		default:
			// If we have a fixed duration, check if time is up
			if !endTime.IsZero() && time.Now().After(endTime) {
				break LOOP
			}
		}

		// Wait for a token from the rate limiter
		if err := t.limiter.Wait(ctx); err != nil {
			// context canceled or limiter error
			break LOOP
		}

		// Acquire concurrency slot
		select {
		case sem <- struct{}{}:
			// proceed
		case <-ctx.Done():
			break LOOP
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			// Simulate a SIP call
			err := sip.SimulateCall(ctx, t.cfg.TargetURI, t.cfg.LocalContact, t.logger)
			t.collector.AddCall(err)
		}()
	}

	// Wait for any in-flight calls to finish
	wg.Wait()

	t.collector.StopTimer()
	return nil
}
