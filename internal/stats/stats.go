package stats

import (
	"sync"
	"time"
)

// Collector collects metrics about the test execution.
type Collector struct {
	mu            sync.Mutex
	totalCalls    int
	totalFailures int
	startTime     time.Time
	endTime       time.Time
}

// NewCollector creates a new stats collector.
func NewCollector() *Collector {
	return &Collector{}
}

// AddCall increments total calls and totalFailures if err != nil.
func (c *Collector) AddCall(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalCalls++
	if err != nil {
		c.totalFailures++
	}
}

// TotalCalls returns the total number of calls attempted.
func (c *Collector) TotalCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.totalCalls
}

// TotalFailures returns the total number of failed calls.
func (c *Collector) TotalFailures() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.totalFailures
}

// StartTimer records the start time.
func (c *Collector) StartTimer() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.startTime = time.Now()
	c.endTime = time.Time{}
}

// StopTimer records the end time.
func (c *Collector) StopTimer() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.endTime = time.Now()
}

// Elapsed returns how long the test ran. If endTime is zero, use the current time.
func (c *Collector) Elapsed() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.startTime.IsZero() {
		return 0
	}
	if c.endTime.IsZero() {
		return time.Since(c.startTime)
	}
	return c.endTime.Sub(c.startTime)
}
