package pipeline

import (
	"strings"

	"github.com/user/logslice/internal/fieldextract"
)

// fieldExtractor wraps the fieldextract package for use in the pipeline.
// It extracts structured fields from a log line and returns them as a
// flat map. The extractor is chosen based on the detected line format.
type fieldExtractor struct {
	jsonEx *fieldextract.Extractor
	kvEx   *fieldextract.Extractor
}

// newFieldExtractor constructs a fieldExtractor that tries JSON extraction
// first, falling back to key=value parsing.
func newFieldExtractor() *fieldExtractor {
	return &fieldExtractor{
		jsonEx: fieldextract.WithJSON(),
		kvEx:   fieldextract.WithKV(),
	}
}

// Extract attempts to pull structured fields from the given log line.
// It returns the extracted fields and a boolean indicating success.
// JSON lines (or lines with a leading timestamp followed by JSON) are
// preferred; key=value lines are tried as a fallback.
func (fe *fieldExtractor) Extract(line string) (map[string]string, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}

	// Try JSON / timestamp+JSON first.
	if fields, ok := fe.jsonEx.Extract(line); ok && len(fields) > 0 {
		return fields, true
	}

	// Fall back to key=value.
	if fields, ok := fe.kvEx.Extract(line); ok && len(fields) > 0 {
		return fields, true
	}

	return nil, false
}

// Level returns the log level found in the extracted fields, checking
// common key names: "level", "lvl", "severity".
func Level(fields map[string]string) string {
	for _, key := range []string{"level", "lvl", "severity"} {
		if v, ok := fields[key]; ok {
			return strings.ToUpper(strings.TrimSpace(v))
		}
	}
	return ""
}
