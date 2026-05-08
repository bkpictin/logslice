// Package fieldextract provides utilities for extracting structured fields
// from unstructured log lines using configurable key=value or JSON patterns.
package fieldextract

import (
	"encoding/json"
	"regexp"
	"strings"
)

// kvPattern matches key=value and key="quoted value" pairs.
var kvPattern = regexp.MustCompile(`(\w+)=("[^"]*"|\S+)`)

// Fields holds the extracted key/value pairs from a log line.
type Fields map[string]string

// Extractor extracts structured fields from a log line.
type Extractor struct {
	tryJSON bool
	tryKV   bool
}

// Option configures an Extractor.
type Option func(*Extractor)

// WithJSON enables JSON field extraction.
func WithJSON() Option {
	return func(e *Extractor) { e.tryJSON = true }
}

// WithKV enables key=value field extraction.
func WithKV() Option {
	return func(e *Extractor) { e.tryKV = true }
}

// New creates an Extractor with the given options.
// By default both JSON and KV extraction are enabled.
func New(opts ...Option) *Extractor {
	e := &Extractor{tryJSON: true, tryKV: true}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Extract attempts to parse structured fields from line.
// JSON objects take priority over key=value scanning when both are enabled.
func (e *Extractor) Extract(line string) Fields {
	line = strings.TrimSpace(line)
	if line == "" {
		return Fields{}
	}

	if e.tryJSON {
		if f := extractJSON(line); len(f) > 0 {
			return f
		}
	}

	if e.tryKV {
		return extractKV(line)
	}

	return Fields{}
}

// extractJSON attempts to decode a JSON object from line.
func extractJSON(line string) Fields {
	// Find first '{' to tolerate leading timestamp prefixes.
	start := strings.IndexByte(line, '{')
	if start < 0 {
		return nil
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line[start:]), &raw); err != nil {
		return nil
	}
	f := make(Fields, len(raw))
	for k, v := range raw {
		switch tv := v.(type) {
		case string:
			f[k] = tv
		default:
			b, _ := json.Marshal(v)
			f[k] = string(b)
		}
	}
	return f
}

// extractKV scans line for key=value pairs.
func extractKV(line string) Fields {
	matches := kvPattern.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return Fields{}
	}
	f := make(Fields, len(matches))
	for _, m := range matches {
		val := strings.Trim(m[2], `"`)
		f[m[1]] = val
	}
	return f
}
