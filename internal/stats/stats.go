package stats

import (
	"sync"
	"time"
)

// Collector accumulates statistics about log processing runs.
type Collector struct {
	mu           sync.Mutex
	TotalLines   int64
	MatchedLines int64
	SkippedLines int64
	BytesRead    int64
	StartTime    time.Time
	EndTime      time.Time
	levelCounts  map[string]int64
}

// New creates a new Collector and records the start time.
func New() *Collector {
	return &Collector{
		StartTime:   time.Now(),
		levelCounts: make(map[string]int64),
	}
}

// AddLine records a processed line. matched indicates whether the line
// passed all filters. level may be empty if not detected.
func (c *Collector) AddLine(matched bool, level string, bytes int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.TotalLines++
	c.BytesRead += bytes
	if matched {
		c.MatchedLines++
		if level != "" {
			c.levelCounts[level]++
		}
	} else {
		c.SkippedLines++
	}
}

// Finish records the end time. Call once processing is complete.
func (c *Collector) Finish() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.EndTime = time.Now()
}

// Duration returns the elapsed time between Start and Finish.
func (c *Collector) Duration() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.EndTime.IsZero() {
		return time.Since(c.StartTime)
	}
	return c.EndTime.Sub(c.StartTime)
}

// LevelCounts returns a copy of the per-level matched line counts.
func (c *Collector) LevelCounts() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[string]int64, len(c.levelCounts))
	for k, v := range c.levelCounts {
		copy[k] = v
	}
	return copy
}
