// Package output provides structured output formatting for logslice.
//
// It supports three output formats:
//
//   - text: human-readable plain text (default)
//     Format: <timestamp> [LEVEL] message
//
//   - json: newline-delimited JSON objects, one per log entry
//     Each object contains timestamp, level, message, and optional fields.
//
//   - csv: comma-separated values with a header row
//     Columns: timestamp, level, message
//
// Example usage:
//
//	f := output.New(os.Stdout, output.FormatJSON)
//	f.Write(output.Entry{
//		Timestamp: time.Now(),
//		Level:     "info",
//		Message:   "server started",
//	})
package output
