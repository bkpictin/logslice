// Package filter provides log entry matching based on configurable criteria
// including time range, log level, and message pattern (regex).
//
// Usage:
//
//	start := time.Now().Add(-1 * time.Hour)
//	opts := filter.Options{
//		StartTime: &start,
//		Level:     "ERROR",
//		Pattern:   "timeout",
//	}
//	f, err := filter.New(opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if f.Match(entry) {
//		// process matching entry
//	}
package filter
