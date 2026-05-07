package pipeline

import "time"

// Config holds all configuration options for a pipeline run.
type Config struct {
	// Files is the list of log file paths to process.
	// If empty, the pipeline reads from stdin.
	Files []string

	// Format controls the output format: "text", "json", or "csv".
	// Defaults to "text" if empty.
	Format string

	// Level filters log lines to those at or above the given level.
	// Recognised values: debug, info, warn, error, fatal.
	Level string

	// Pattern is an optional regular-expression pattern.
	// Only lines whose message matches the pattern are emitted.
	Pattern string

	// Since discards lines with a timestamp before this time.
	// Zero value means no lower bound.
	Since time.Time

	// Until discards lines with a timestamp after this time.
	// Zero value means no upper bound.
	Until time.Time

	// MaxLines stops processing after this many matched lines.
	// Zero means unlimited.
	MaxLines int
}

// format returns the effective output format, defaulting to "text".
func (c Config) format() string {
	if c.Format == "" {
		return "text"
	}
	return c.Format
}
