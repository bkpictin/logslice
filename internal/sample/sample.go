// Package sample provides log line sampling to reduce output volume
// while preserving a statistically representative subset of lines.
package sample

import (
	"sync/atomic"
)

// Sampler decides whether a given log line should be kept based on a
// 1-in-N sampling strategy. A rate of 1 keeps every line; a rate of N
// keeps approximately 1/N of all lines.
type Sampler struct {
	rate    uint64
	counter atomic.Uint64
	skipped atomic.Uint64
}

// New returns a Sampler that retains one line for every `rate` lines
// seen. A rate <= 1 disables sampling (all lines are kept).
func New(rate uint64) *Sampler {
	if rate < 1 {
		rate = 1
	}
	return &Sampler{rate: rate}
}

// Keep returns true if the current line should be emitted.
// It is safe to call concurrently.
func (s *Sampler) Keep() bool {
	n := s.counter.Add(1)
	if n%s.rate == 1 {
		return true
	}
	s.skipped.Add(1)
	return false
}

// Skipped returns the total number of lines that were dropped by the
// sampler since creation or the last Reset.
func (s *Sampler) Skipped() uint64 {
	return s.skipped.Load()
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() uint64 {
	return s.rate
}

// Enabled reports whether sampling is active (rate > 1).
func (s *Sampler) Enabled() bool {
	return s.rate > 1
}

// Reset clears internal counters, restarting the sampling window.
func (s *Sampler) Reset() {
	s.counter.Store(0)
	s.skipped.Store(0)
}
