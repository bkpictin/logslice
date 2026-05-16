package masklevel

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/levelmap"
)

// Masker replaces the log level token in a line with a fixed replacement string.
// This is useful when forwarding logs to systems that expect a normalised level
// field, or when you want to redact severity information from output.
type Masker struct {
	minLevel  levelmap.Level
	mask      string
	enabled   bool
}

// Option configures a Masker.
type Option func(*Masker)

// WithMask overrides the default replacement string ("[LEVEL]").
func WithMask(s string) Option {
	return func(m *Masker) { m.mask = s }
}

// New returns a Masker that rewrites any level token at or above minLevel.
// If minLevel is unknown the Masker is disabled and Line is a no-op.
func New(minLevel string, opts ...Option) (*Masker, error) {
	lvl, ok := levelmap.Parse(minLevel)
	if !ok && minLevel != "" {
		return nil, fmt.Errorf("masklevel: unknown level %q", minLevel)
	}
	m := &Masker{
		minLevel: lvl,
		mask:     "[LEVEL]",
		enabled:  ok,
	}
	for _, o := range opts {
		o(m)
	}
	return m, nil
}

// Enabled reports whether the masker will modify lines.
func (m *Masker) Enabled() bool { return m.enabled }

// Line scans line for a known level token. If the token is at or above the
// configured minimum level it is replaced with the mask string.
func (m *Masker) Line(line string) string {
	if !m.enabled {
		return line
	}
	for _, name := range levelmap.Names() {
		lvl, _ := levelmap.Parse(name)
		if !levelmap.AtLeast(lvl, m.minLevel) {
			continue
		}
		if idx := strings.Index(strings.ToUpper(line), name); idx >= 0 {
			// Replace the first occurrence preserving surrounding text.
			upper := strings.ToUpper(line)
			return line[:idx] + m.mask + line[idx+len(name):]
			_ = upper
		}
	}
	return line
}
