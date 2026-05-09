// Package burst provides a sliding-window burst detector that flags when
// the number of matched log lines exceeds a threshold within a time window.
package burst

import (
	"sync"
	"time"
)

// Detector tracks log line timestamps in a sliding window and reports
// whether the current burst rate exceeds the configured threshold.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	times     []time.Time
	triggered int
}

// New creates a Detector that fires when more than threshold events
// occur within the given window duration. A zero or negative threshold
// disables burst detection (Observe always returns false).
func New(threshold int, window time.Duration) *Detector {
	return &Detector{
		threshold: threshold,
		window:    window,
	}
}

// Observe records a new event at the given timestamp and returns true
// if the burst threshold has been exceeded within the sliding window.
func (d *Detector) Observe(at time.Time) bool {
	if d.threshold <= 0 {
		return false
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := at.Add(-d.window)
	d.evict(cutoff)
	d.times = append(d.times, at)

	if len(d.times) > d.threshold {
		d.triggered++
		return true
	}
	return false
}

// TriggeredCount returns the total number of times the burst threshold
// has been exceeded since the detector was created or last reset.
func (d *Detector) TriggeredCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.triggered
}

// Reset clears all recorded events and resets the triggered counter.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.times = d.times[:0]
	d.triggered = 0
}

// evict removes timestamps older than cutoff. Must be called with mu held.
func (d *Detector) evict(cutoff time.Time) {
	i := 0
	for i < len(d.times) && d.times[i].Before(cutoff) {
		i++
	}
	d.times = d.times[i:]
}
