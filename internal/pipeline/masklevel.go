package pipeline

import (
	"fmt"

	"github.com/logslice/logslice/internal/masklevel"
)

// newMaskLevel constructs a masklevel.Masker from pipeline config.
// Returns nil when no mask level is configured so callers can skip the step.
func newMaskLevel(cfg Config) (*masklevel.Masker, error) {
	if cfg.MaskLevel == "" {
		return nil, nil
	}
	opts := []masklevel.Option{}
	if cfg.MaskLevelToken != "" {
		opts = append(opts, masklevel.WithMask(cfg.MaskLevelToken))
	}
	m, err := masklevel.New(cfg.MaskLevel, opts...)
	if err != nil {
		return nil, fmt.Errorf("pipeline: masklevel: %w", err)
	}
	return m, nil
}

// applyMaskLevel rewrites the level token in line when the masker is active.
// A nil masker is a no-op.
func applyMaskLevel(m *masklevel.Masker, line string) string {
	if m == nil || !m.Enabled() {
		return line
	}
	return m.Line(line)
}
