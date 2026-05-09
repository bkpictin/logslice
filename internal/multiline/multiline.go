// Package multiline provides support for collapsing multi-line log entries
// (e.g. Java stack traces, Python tracebacks) into a single logical record.
package multiline

import (
	"regexp"
	"strings"
)

// Folder accumulates lines that belong to a single logical log entry.
// A new entry begins when a line matches the Start pattern. Lines that
// do NOT match Start are treated as continuation lines and appended to
// the current entry.
type Folder struct {
	start     *regexp.Regexp
	current   strings.Builder
	hasEntry  bool
	Separator string // separator inserted between continuation lines (default "\n")
}

// Option configures a Folder.
type Option func(*Folder)

// WithSeparator overrides the string used to join continuation lines.
func WithSeparator(sep string) Option {
	return func(f *Folder) { f.Separator = sep }
}

// New creates a Folder whose entries begin with lines matching startPattern.
// Returns an error if startPattern is not a valid regular expression.
func New(startPattern string, opts ...Option) (*Folder, error) {
	re, err := regexp.Compile(startPattern)
	if err != nil {
		return nil, err
	}
	f := &Folder{
		start:     re,
		Separator: "\n",
	}
	for _, o := range opts {
		o(f)
	}
	return f, nil
}

// Feed adds a raw line to the folder. If the line starts a new entry the
// previous entry is returned (flushed); otherwise nil is returned.
func (f *Folder) Feed(line string) (flushed string, ok bool) {
	if f.start.MatchString(line) {
		if f.hasEntry {
			prev := f.current.String()
			f.current.Reset()
			f.current.WriteString(line)
			return prev, true
		}
		f.hasEntry = true
		f.current.WriteString(line)
		return "", false
	}
	// continuation line
	if f.hasEntry {
		f.current.WriteString(f.Separator)
		f.current.WriteString(line)
	}
	return "", false
}

// Flush returns any buffered entry and resets the folder.
// Call this after the input stream is exhausted.
func (f *Folder) Flush() (string, bool) {
	if !f.hasEntry {
		return "", false
	}
	result := f.current.String()
	f.current.Reset()
	f.hasEntry = false
	return result, true
}

// Reset discards any buffered state without emitting.
func (f *Folder) Reset() {
	f.current.Reset()
	f.hasEntry = false
}
