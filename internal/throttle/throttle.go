// Package throttle provides a line-based throughput limiter that drops
// lines once a configurable lines-per-second ceiling is exceeded within
// a sliding time window.
package throttle

import (
	"sync"
	"time"
)

// Throttler tracks line throughput and signals when the rate limit is
// exceeded. A zero-value Throttler (or one constructed with maxLPS <= 0)
// is disabled and never drops lines.
type Throttler struct {
	maxLPS  int
	window  time.Duration
	mu      sync.Mutex
	buckets []time.Time
	dropped int64
}

// New returns a Throttler that allows at most maxLPS lines per second.
// Pass maxLPS <= 0 to disable throttling entirely.
func New(maxLPS int) *Throttler {
	return &Throttler{
		maxLPS: maxLPS,
		window: time.Second,
	}
}

// Enabled reports whether throttling is active.
func (t *Throttler) Enabled() bool {
	return t.maxLPS > 0
}

// Allow returns true if the current line should be passed through, or
// false if it should be dropped because the rate limit has been reached.
// Allow is safe for concurrent use.
func (t *Throttler) Allow() bool {
	if !t.Enabled() {
		return true
	}
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := now.Add(-t.window)
	valid := t.buckets[:0]
	for _, ts := range t.buckets {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	t.buckets = valid

	if len(t.buckets) >= t.maxLPS {
		t.dropped++
		return false
	}
	t.buckets = append(t.buckets, now)
	return true
}

// Dropped returns the total number of lines dropped since creation or
// the last Reset.
func (t *Throttler) Dropped() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.dropped
}

// Reset clears internal state, including the dropped counter.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buckets = t.buckets[:0]
	t.dropped = 0
}
