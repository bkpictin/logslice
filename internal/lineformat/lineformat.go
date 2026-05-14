// Package lineformat normalises raw log lines before further processing.
// It trims trailing whitespace and carriage returns, optionally strips ANSI
// escape sequences, and enforces a maximum byte length on the raw input so
// that downstream components never receive pathologically long lines.
package lineformat

import (
	"regexp"
	"strings"
)

// ansiEscape matches ANSI/VT100 CSI escape sequences.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// Options controls the behaviour of a Formatter.
type Options struct {
	// StripANSI removes ANSI colour/control codes from each line.
	StripANSI bool
	// MaxBytes caps the input line at this many bytes before any other
	// processing. Zero means no limit.
	MaxBytes int
}

// Formatter normalises log lines.
type Formatter struct {
	opts Options
}

// New returns a Formatter configured with opts.
func New(opts Options) *Formatter {
	return &Formatter{opts: opts}
}

// NewDefault returns a Formatter with ANSI stripping enabled and no byte cap.
func NewDefault() *Formatter {
	return New(Options{StripANSI: true})
}

// Line normalises a single log line and returns the result.
// The original string is never mutated.
func (f *Formatter) Line(line string) string {
	// Cap raw bytes first so we don't waste work on huge lines.
	if f.opts.MaxBytes > 0 && len(line) > f.opts.MaxBytes {
		line = line[:f.opts.MaxBytes]
	}

	// Trim trailing CR/LF and spaces.
	line = strings.TrimRight(line, "\r\n ")

	if f.opts.StripANSI {
		line = ansiEscape.ReplaceAllString(line, "")
	}

	return line
}

// Lines applies Line to every element of lines and returns a new slice.
func (f *Formatter) Lines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = f.Line(l)
	}
	return out
}

// Enabled reports whether the formatter will perform any transformation.
func (f *Formatter) Enabled() bool {
	return f.opts.StripANSI || f.opts.MaxBytes > 0
}
