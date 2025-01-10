// File: internal/stats/stats_test.go
package stats_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/copyleftdev/sipload/internal/stats"
)

func TestCollector_Basic(t *testing.T) {
	c := stats.NewCollector()

	c.StartTimer()
	time.Sleep(50 * time.Millisecond)   // simulate some passing time
	c.AddCall(nil)                      // success
	c.AddCall(errors.New("mock error")) // fail
	c.StopTimer()

	assert.Equal(t, 2, c.TotalCalls(), "Total calls should be 2")
	assert.Equal(t, 1, c.TotalFailures(), "Should have 1 failure")

	elapsed := c.Elapsed()
	assert.GreaterOrEqual(t, int64(elapsed), int64(50*time.Millisecond))
}
