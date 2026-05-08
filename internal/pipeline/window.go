package pipeline

import (
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/window"
)

// WindowConfig controls the context-window extraction step in the pipeline.
type WindowConfig struct {
	// Enabled turns context-window extraction on.
	Enabled bool
	// Before is how far back in time to capture lines before an anchor match.
	Before time.Duration
	// After is how far forward in time to capture lines after an anchor match.
	After time.Duration
}

// newWindowFilter builds a window.Window from a WindowConfig, or returns nil
// when the feature is disabled.
func newWindowFilter(cfg WindowConfig) *window.Window {
	if !cfg.Enabled {
		return nil
	}
	return window.New(cfg.Before, cfg.After)
}

// applyWindow feeds a line through the context window and returns the lines
// that should be forwarded downstream. When w is nil every line passes through
// unchanged (feature is disabled).
//
// The entry timestamp is used as the window clock; if the entry has no parsed
// time the current wall clock is used as a fallback so the window still works
// with unstructured logs.
func applyWindow(w *window.Window, line string, entry *filter.Entry, match bool) []string {
	if w == nil {
		if match {
			return []string{line}
		}
		return nil
	}

	at := entry.Time
	if at.IsZero() {
		at = time.Now()
	}
	return w.Feed(line, at, match)
}
