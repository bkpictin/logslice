// Package window implements a context-window extractor for log lines.
//
// It buffers recent log entries and, when an anchor line is matched, emits
// all lines within a configurable duration before and after that anchor.
// This is useful for surfacing surrounding context around errors or events
// without requiring line-number counting.
//
// Basic usage:
//
//	w := window.New(5*time.Second, 10*time.Second)
//	for _, entry := range logEntries {
//		lines := w.Feed(entry.Line, entry.Timestamp, isError(entry))
//		for _, l := range lines {
//			fmt.Println(l)
//		}
//	}
package window
