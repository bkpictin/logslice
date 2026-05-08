// Package fieldextract provides utilities for extracting structured fields
// from raw log lines. It supports JSON log lines, key=value pairs, and
// hybrid formats where a timestamp or prefix precedes a JSON payload.
//
// Extracted fields are returned as a map[string]string and can be used
// downstream for filtering, formatting, or enrichment.
//
// Supported formats:
//
//	{"level":"info","msg":"started"}          // pure JSON
//	2024-01-01T00:00:00Z {"level":"info"}     // timestamp + JSON
//	level=info msg=started                    // key=value pairs
package fieldextract
