package pipeline

import (
	"github.com/logslice/logslice/internal/levelmap"
	"github.com/logslice/logslice/internal/severity"
)

// newSeverityFilter constructs a severity.Filter from the pipeline Config.
// Returns nil when no minimum level is configured so callers can skip the
// allocation entirely.
func newSeverityFilter(cfg Config) *severity.Filter {
	if cfg.MinLevel == "" {
		return nil
	}
	min := levelmap.Parse(cfg.MinLevel)
	if min == levelmap.Unknown {
		return nil
	}
	return severity.New(min)
}

// applySeverity returns false when line should be dropped according to f.
// level is the extracted log level string for the current line.
// If f is nil every line is allowed.
func applySeverity(f *severity.Filter, level string) bool {
	if f == nil {
		return true
	}
	return f.Allow(level)
}
