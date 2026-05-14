package pipeline

import (
	"github.com/user/logslice/internal/dedupe"
	"github.com/user/logslice/internal/filter"
)

// newDeduplicator constructs a dedupe.Deduper from the pipeline Config.
// Returns nil when deduplication is disabled (WindowSize == 0).
func newDeduplicator(cfg Config) (*dedupe.Deduper, error) {
	if cfg.DedupeWindow <= 0 {
		return nil, nil
	}
	return dedupe.New(cfg.DedupeWindow), nil
}

// applyDedupe filters the entry through the deduplicator when one is
// configured. It returns (entry, true) when the line should be kept and
// (entry, false) when it is a duplicate that should be dropped.
func applyDedupe(d *dedupe.Deduper, entry filter.Entry) (filter.Entry, bool) {
	if d == nil {
		return entry, true
	}
	if d.Seen(entry.Line) {
		return entry, false
	}
	return entry, true
}
