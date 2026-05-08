package pipeline

import (
	"fmt"

	"github.com/yourorg/logslice/internal/redact"
)

// newRedactor builds a *redact.Redactor from the pipeline Config.
// If no redact patterns are configured it returns a disabled no-op
// redactor so callers never need to nil-check.
func newRedactor(cfg Config) (*redact.Redactor, error) {
	var opts []redact.Option
	if cfg.RedactPlaceholder != "" {
		opts = append(opts, redact.WithPlaceholder(cfg.RedactPlaceholder))
	}
	r, err := redact.New(cfg.RedactPatterns, opts...)
	if err != nil {
		return nil, fmt.Errorf("redact: %w", err)
	}
	return r, nil
}

// applyRedaction rewrites line using the supplied Redactor.
// When the Redactor is disabled the original line is returned unchanged.
func applyRedaction(r *redact.Redactor, line string) string {
	if r == nil || !r.Enabled() {
		return line
	}
	return r.Line(line)
}
