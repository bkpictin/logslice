package pipeline

import (
	"github.com/yourorg/logslice/internal/truncate"
)

// truncateConfig holds line-length limiting options surfaced via Config.
type truncateConfig struct {
	// MaxLineLen caps each output line at this many bytes.
	// A value of 0 disables truncation.
	MaxLineLen int
	// Suffix is appended to truncated lines. Defaults to "..." when empty
	// and MaxLineLen > 0.
	Suffix string
}

// newTruncator builds a Truncator from the pipeline Config, applying
// sensible defaults. Returns a no-op Truncator when truncation is disabled.
func newTruncator(cfg Config) *truncate.Truncator {
	if cfg.MaxLineLen <= 0 {
		return truncate.New(0, "")
	}
	suffix := cfg.Suffix
	if suffix == "" {
		suffix = truncate.DefaultSuffix
	}
	return truncate.New(cfg.MaxLineLen, suffix)
}

// applyTruncation runs tr.Line over line, returning the (possibly shortened)
// result. Inlined here so callers in pipeline.go stay readable.
func applyTruncation(tr *truncate.Truncator, line string) string {
	return tr.Line(line)
}
