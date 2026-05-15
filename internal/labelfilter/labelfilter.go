package labelfilter

import "strings"

// Filter matches log lines that contain all required key=value labels.
// Labels are matched as substrings of the form "key=value" within a line.
type Filter struct {
	labels   map[string]string
	enabled  bool
	matched  int
	rejected int
}

// New creates a Filter from a slice of "key=value" strings.
// If pairs is empty the filter is disabled and every line passes.
func New(pairs []string) (*Filter, error) {
	labels := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok || strings.TrimSpace(k) == "" {
			return nil, &ParseError{Raw: p}
		}
		labels[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return &Filter{labels: labels, enabled: len(labels) > 0}, nil
}

// Enabled reports whether the filter has any labels to match.
func (f *Filter) Enabled() bool { return f.enabled }

// Match returns true when the line contains every configured label pair.
// If the filter is disabled it always returns true.
func (f *Filter) Match(line string) bool {
	if !f.enabled {
		return true
	}
	for k, v := range f.labels {
		needle := k + "=" + v
		if !strings.Contains(line, needle) {
			f.rejected++
			return false
		}
	}
	f.matched++
	return true
}

// Matched returns the number of lines that passed the filter.
func (f *Filter) Matched() int { return f.matched }

// Rejected returns the number of lines that were dropped.
func (f *Filter) Rejected() int { return f.rejected }

// Reset clears counters without changing the label set.
func (f *Filter) Reset() {
	f.matched = 0
	f.rejected = 0
}

// ParseError is returned when a label pair cannot be parsed.
type ParseError struct{ Raw string }

func (e *ParseError) Error() string {
	return "labelfilter: invalid label pair: " + e.Raw
}
