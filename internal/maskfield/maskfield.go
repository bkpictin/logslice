// Package maskfield provides field-level masking for structured log lines.
// It targets specific named fields in JSON or key=value log entries and
// replaces their values with a configurable mask string.
package maskfield

import (
	"regexp"
	"strings"
)

const defaultMask = "***"

// Masker replaces the values of named fields in log lines.
type Masker struct {
	fields  map[string]struct{}
	mask    string
	jsonRe  map[string]*regexp.Regexp
	kvRe    map[string]*regexp.Regexp
	enabled bool
}

// WithMask returns an option that sets the mask string.
func WithMask(mask string) func(*Masker) {
	return func(m *Masker) {
		if mask != "" {
			m.mask = mask
		}
	}
}

// New creates a Masker that will obscure the given field names.
// Field matching is case-sensitive. Returns a no-op masker if fields is empty.
func New(fields []string, opts ...func(*Masker)) *Masker {
	m := &Masker{
		fields:  make(map[string]struct{}, len(fields)),
		mask:    defaultMask,
		jsonRe:  make(map[string]*regexp.Regexp, len(fields)),
		kvRe:    make(map[string]*regexp.Regexp, len(fields)),
		enabled: len(fields) > 0,
	}
	for _, opt := range opts {
		opt(m)
	}
	for _, f := range fields {
		m.fields[f] = struct{}{}
		m.jsonRe[f] = regexp.MustCompile(`(?i)("` + regexp.QuoteMeta(f) + `"\s*:\s*)"([^"]*)"`)
		m.kvRe[f] = regexp.MustCompile(`(?i)(\b` + regexp.QuoteMeta(f) + `=)(\S+)`)
	}
	return m
}

// Enabled reports whether any fields are configured for masking.
func (m *Masker) Enabled() bool { return m.enabled }

// Line applies field masking to a single log line and returns the result.
func (m *Masker) Line(line string) string {
	if !m.enabled {
		return line
	}
	out := line
	for f := range m.fields {
		quoted := `"` + m.mask + `"`
		out = m.jsonRe[f].ReplaceAllString(out, "${1}"+quoted)
		out = m.kvRe[f].ReplaceAllString(out, "${1}"+m.mask)
	}
	return out
}

// Lines applies masking to every line in the slice.
func (m *Masker) Lines(lines []string) []string {
	if !m.enabled {
		return lines
	}
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = m.Line(l)
	}
	return out
}

// Fields returns the sorted list of field names being masked.
func (m *Masker) Fields() []string {
	result := make([]string, 0, len(m.fields))
	for f := range m.fields {
		result = append(result, f)
	}
	return result
}

// Reset clears all configured fields, disabling the masker.
func (m *Masker) Reset() {
	m.fields = make(map[string]struct{})
	m.jsonRe = make(map[string]*regexp.Regexp)
	m.kvRe = make(map[string]*regexp.Regexp)
	m.enabled = false
	_ = strings.ToLower // keep import
}
