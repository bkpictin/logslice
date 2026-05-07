package highlight

import (
	"regexp"
	"strings"
)

// ANSI colour codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// LevelColor returns the ANSI colour associated with a log level string.
func LevelColor(level string) string {
	switch strings.ToUpper(level) {
	case "ERROR", "FATAL", "CRITICAL":
		return Red
	case "WARN", "WARNING":
		return Yellow
	case "INFO":
		return Green
	case "DEBUG", "TRACE":
		return Cyan
	default:
		return ""
	}
}

// Highlighter applies ANSI colours to log lines.
type Highlighter struct {
	pattern *regexp.Regexp
	enabled bool
}

// New creates a Highlighter. When enabled is false all methods are no-ops.
// pattern may be nil; when set, matches within a line are bolded.
func New(enabled bool, pattern *regexp.Regexp) *Highlighter {
	return &Highlighter{enabled: enabled, pattern: pattern}
}

// Line wraps the entire line with the colour for level, then bolds any
// pattern matches inside it.
func (h *Highlighter) Line(line, level string) string {
	if !h.enabled {
		return line
	}
	if h.pattern != nil {
		line = h.pattern.ReplaceAllStringFunc(line, func(m string) string {
			return Bold + m + Reset
		})
	}
	color := LevelColor(level)
	if color == "" {
		return line
	}
	return color + line + Reset
}

// Word wraps a single token with the colour for level.
func (h *Highlighter) Word(word, level string) string {
	if !h.enabled {
		return word
	}
	color := LevelColor(level)
	if color == "" {
		return word
	}
	return color + word + Reset
}
