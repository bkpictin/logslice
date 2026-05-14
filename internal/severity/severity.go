// Package severity provides min-severity filtering for log entries,
// allowing callers to suppress lines below a configured threshold.
package severity

import "github.com/logslice/logslice/internal/levelmap"

// Filter drops log entries whose level is below the configured minimum.
type Filter struct {
	min     levelmap.Level
	enabled bool
	dropped int
}

// New returns a Filter that passes only entries at or above min.
// If min is levelmap.Unknown the filter is disabled and all lines pass.
func New(min levelmap.Level) *Filter {
	return &Filter{
		min:     min,
		enabled: min != levelmap.Unknown,
	}
}

// Allow returns true when the line should be forwarded downstream.
// level is the parsed level of the current log line; an empty string is
// treated as Unknown and passes when the filter is enabled only if the
// minimum is also Unknown.
func (f *Filter) Allow(level string) bool {
	if !f.enabled {
		return true
	}
	parsed := levelmap.Parse(level)
	if !levelmap.AtLeast(parsed, f.min) {
		f.dropped++
		return false
	}
	return true
}

// Dropped returns the number of lines suppressed since creation or last Reset.
func (f *Filter) Dropped() int { return f.dropped }

// Enabled reports whether the filter is active.
func (f *Filter) Enabled() bool { return f.enabled }

// Reset clears the dropped counter.
func (f *Filter) Reset() { f.dropped = 0 }
