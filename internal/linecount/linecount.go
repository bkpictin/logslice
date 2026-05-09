// Package linecount provides a simple rolling line counter that tracks
// total lines seen and emits a progress callback every N lines.
package linecount

import "sync/atomic"

// Counter tracks lines processed and fires a callback at a configurable
// interval. It is safe for concurrent use.
type Counter struct {
	total    atomic.Int64
	interval int64
	cb       func(n int64)
}

// New creates a Counter that calls cb every interval lines.
// If interval is zero or negative, the callback is never called.
func New(interval int, cb func(n int64)) *Counter {
	if interval < 0 {
		interval = 0
	}
	return &Counter{
		interval: int64(interval),
		cb:       cb,
	}
}

// Add increments the counter by delta and fires the callback when the
// cumulative total crosses a multiple of the configured interval.
func (c *Counter) Add(delta int64) {
	if delta <= 0 {
		return
	}
	prev := c.total.Load()
	next := c.total.Add(delta)

	if c.interval <= 0 || c.cb == nil {
		return
	}

	// Fire once per interval boundary crossed.
	prevTick := prev / c.interval
	nextTick := (next - 1) / c.interval
	for i := prevTick + 1; i <= nextTick; i++ {
		c.cb(i * c.interval)
	}
}

// Inc increments the counter by one.
func (c *Counter) Inc() { c.Add(1) }

// Total returns the current cumulative line count.
func (c *Counter) Total() int64 { return c.total.Load() }

// Reset sets the counter back to zero.
func (c *Counter) Reset() { c.total.Store(0) }
