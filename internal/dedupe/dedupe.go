// Package dedupe provides line deduplication for log streams.
// It tracks recently seen lines and suppresses consecutive or
// near-consecutive duplicates, optionally reporting a suppression count.
package dedupe

import "sync"

// Deduper filters duplicate log lines within a sliding window.
type Deduper struct {
	mu      sync.Mutex
	window  int
	recent  []string
	pos     int
	count   int
	Skipped int
}

// New returns a Deduper that remembers the last windowSize unique lines.
// A windowSize of 1 suppresses only immediately repeated lines.
func New(windowSize int) *Deduper {
	if windowSize < 1 {
		windowSize = 1
	}
	return &Deduper{
		window: windowSize,
		recent: make([]string, windowSize),
	}
}

// IsDuplicate returns true if line was seen within the current window.
// If it is a duplicate, Skipped is incremented and the line is not
// added to the window again. If it is new, it is recorded.
func (d *Deduper) IsDuplicate(line string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, s := range d.recent {
		if s == line {
			d.Skipped++
			return true
		}
	}

	d.recent[d.pos%d.window] = line
	d.pos++
	d.count++
	return false
}

// Reset clears the window and counters.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.recent = make([]string, d.window)
	d.pos = 0
	d.count = 0
	d.Skipped = 0
}

// Seen returns the number of unique lines recorded so far.
func (d *Deduper) Seen() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.count
}
