package redact

import (
	"regexp"
	"strings"
)

// Redactor replaces sensitive patterns in log lines with a placeholder.
type Redactor struct {
	patterns []*regexp.Regexp
	placeholder string
	enabled bool
}

// Option configures a Redactor.
type Option func(*Redactor)

// WithPlaceholder sets the replacement string (default: "[REDACTED]").
func WithPlaceholder(s string) Option {
	return func(r *Redactor) {
		if s != "" {
			r.placeholder = s
		}
	}
}

// New creates a Redactor from a list of regex pattern strings.
// If no patterns are provided the Redactor is disabled.
func New(patterns []string, opts ...Option) (*Redactor, error) {
	r := &Redactor{
		placeholder: "[REDACTED]",
	}
	for _, o := range opts {
		o(r)
	}
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		r.patterns = append(r.patterns, re)
	}
	r.enabled = len(r.patterns) > 0
	return r, nil
}

// Enabled reports whether any patterns are configured.
func (r *Redactor) Enabled() bool {
	return r.enabled
}

// Line applies all redaction patterns to a single log line.
func (r *Redactor) Line(line string) string {
	if !r.enabled {
		return line
	}
	for _, re := range r.patterns {
		line = re.ReplaceAllString(line, r.placeholder)
	}
	return line
}

// Lines applies redaction to a slice of lines in place.
func (r *Redactor) Lines(lines []string) []string {
	if !r.enabled {
		return lines
	}
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = r.Line(l)
	}
	return out
}
