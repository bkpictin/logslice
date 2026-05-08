package truncate

import "strings"

// Truncator trims long log lines to a maximum byte length,
// optionally appending a suffix such as "..." to indicate truncation.
type Truncator struct {
	maxLen int
	suffix string
}

// New returns a Truncator that limits lines to maxLen bytes.
// If maxLen <= 0 truncation is disabled (lines pass through unchanged).
// suffix is appended when a line is shortened; pass "" for no marker.
func New(maxLen int, suffix string) *Truncator {
	return &Truncator{maxLen: maxLen, suffix: suffix}
}

// Line returns the (possibly truncated) version of s.
func (t *Truncator) Line(s string) string {
	if t.maxLen <= 0 || len(s) <= t.maxLen {
		return s
	}
	cutAt := t.maxLen
	if suf := t.suffix; suf != "" {
		cutAt = t.maxLen - len(suf)
		if cutAt < 0 {
			cutAt = 0
		}
		return s[:cutAt] + suf
	}
	return s[:cutAt]
}

// Lines applies Line to every element of lines, returning a new slice.
func (t *Truncator) Lines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = t.Line(l)
	}
	return out
}

// Enabled reports whether truncation is active.
func (t *Truncator) Enabled() bool { return t.maxLen > 0 }

// DefaultSuffix is the suffix used by NewDefault.
const DefaultSuffix = "..."

// NewDefault returns a Truncator with the standard "..." suffix.
func NewDefault(maxLen int) *Truncator {
	return New(maxLen, DefaultSuffix)
}

// TrimFields shortens each whitespace-separated field of s independently,
// then reassembles the line with single spaces.
func TrimFields(s string, maxFieldLen int) string {
	if maxFieldLen <= 0 {
		return s
	}
	fields := strings.Fields(s)
	for i, f := range fields {
		if len(f) > maxFieldLen {
			fields[i] = f[:maxFieldLen] + DefaultSuffix
		}
	}
	return strings.Join(fields, " ")
}
