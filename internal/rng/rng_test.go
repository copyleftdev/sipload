// File: internal/rng/rng_test.go
package rng_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/copyleftdev/sipload/internal/rng"
)

// TestSimpleRNG verifies that our RNG returns values in the expected range.
func TestSimpleRNG(t *testing.T) {
	maxVal := 10
	for i := 0; i < 100; i++ {
		val := rng.Intn(maxVal)
		assert.GreaterOrEqual(t, val, 0)
		assert.Less(t, val, maxVal)
	}
}
