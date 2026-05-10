// Package levelmap provides normalisation of log level strings to a
// canonical set of severity values used across logslice.
package levelmap

import "strings"

// Level represents a canonical log severity.
type Level int

const (
	Unknown Level = iota
	Trace
	Debug
	Info
	Warn
	Error
	Fatal
)

// String returns the canonical string representation of a Level.
func (l Level) String() string {
	switch l {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// aliases maps common variant spellings to a canonical Level.
var aliases = map[string]Level{
	"trace":   Trace,
	"trc":     Trace,
	"debug":   Debug,
	"dbg":     Debug,
	"info":    Info,
	"inf":     Info,
	"warn":    Warn,
	"warning": Warn,
	"wrn":     Warn,
	"error":   Error,
	"err":     Error,
	"fatal":   Fatal,
	"crit":    Fatal,
	"critical":Fatal,
	"panic":   Fatal,
}

// Parse converts a raw level string (case-insensitive) to a canonical Level.
// Unknown or empty strings return Unknown.
func Parse(raw string) Level {
	if raw == "" {
		return Unknown
	}
	if l, ok := aliases[strings.ToLower(strings.TrimSpace(raw))]; ok {
		return l
	}
	return Unknown
}

// AtLeast reports whether l is at least as severe as min.
func AtLeast(l, min Level) bool {
	return l >= min
}
