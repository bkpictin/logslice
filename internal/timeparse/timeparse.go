// Package timeparse provides flexible timestamp parsing for log lines.
// It tries a set of common log timestamp formats in order and returns
// the first successful parse result along with the remaining line text.
package timeparse

import (
	"strings"
	"time"
)

// Result holds the parsed time and the remainder of the line after the timestamp.
type Result struct {
	Time      time.Time
	Remainder string
	Format    string
}

// knownFormats lists timestamp layouts tried in order, from most to least specific.
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan  2 15:04:05",
	"Jan 02 15:04:05",
}

// Parse attempts to extract a timestamp from the beginning of line.
// It tries each known format against progressively longer prefixes of the line.
// Returns (Result, true) on success or (zero Result, false) if no format matched.
func Parse(line string) (Result, bool) {
	for _, layout := range knownFormats {
		prefixLen := len(layout) + 10 // generous upper bound
		if prefixLen > len(line) {
			prefixLen = len(line)
		}
		candidate := line[:prefixLen]
		// Trim trailing partial tokens that may confuse the parser.
		for i := prefixLen; i >= len(layout)-2; i-- {
			t, err := time.Parse(layout, candidate[:i])
			if err == nil {
				remainder := strings.TrimLeft(line[i:], " \t")
				return Result{Time: t, Remainder: remainder, Format: layout}, true
			}
		}
	}
	return Result{}, false
}

// Formats returns the ordered list of timestamp layouts that Parse will attempt.
func Formats() []string {
	out := make([]string, len(knownFormats))
	copy(out, knownFormats)
	return out
}
