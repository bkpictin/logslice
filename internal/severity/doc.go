// Package severity implements a minimum-severity gate for log pipelines.
//
// Usage:
//
//	f := severity.New(levelmap.Parse("warn"))
//	if f.Allow(entry.Level) {
//		// forward the line
//	}
//
// The filter is a no-op when constructed with levelmap.Unknown, ensuring
// that pipelines without a --min-level flag incur no overhead.
package severity
