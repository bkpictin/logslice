package pipeline

import (
	"fmt"

	"github.com/user/logslice/internal/labelfilter"
)

// newLabelFilter constructs a labelfilter.Filter from the pipeline Config.
// It returns nil when no label pairs are configured so callers can skip the
// allocation entirely.
func newLabelFilter(cfg Config) (*labelfilter.Filter, error) {
	if len(cfg.Labels) == 0 {
		return nil, nil
	}
	f, err := labelfilter.New(cfg.Labels)
	if err != nil {
		return nil, fmt.Errorf("pipeline: label filter: %w", err)
	}
	return f, nil
}

// applyLabelFilter returns false when the line should be dropped by the label
// filter.  A nil filter always passes.
func applyLabelFilter(f *labelfilter.Filter, line string) bool {
	if f == nil {
		return true
	}
	return f.Match(line)
}
