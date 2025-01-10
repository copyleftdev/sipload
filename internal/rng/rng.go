package rng

import (
	"time"
)

// A very simple linear congruential generator (LCG) for demonstration.
// Not cryptographically secure.
type SimpleRNG struct {
	seed uint64
}

func NewSimpleRNG(seed int64) *SimpleRNG {
	return &SimpleRNG{seed: uint64(seed)}
}

func (r *SimpleRNG) Int63n(n int64) int64 {
	r.seed = (r.seed*6364136223846793005 + 1) % (1 << 63)
	return int64(r.seed % uint64(n))
}

// A global RNG we can use
var globalRNG = NewSimpleRNG(time.Now().UnixNano())

// Intn returns an int in [0, n).
func Intn(n int) int {
	return int(globalRNG.Int63n(int64(n)))
}
