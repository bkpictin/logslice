// Package stats provides a thread-safe collector for log processing
// statistics. It tracks total, matched, and skipped line counts, bytes
// read, elapsed time, and per-level breakdown of matched lines.
//
// Typical usage:
//
//	col := stats.New()
//	for _, line := range lines {
//		col.AddLine(matched, level, int64(len(line)))
//	}
//	col.Finish()
//	fmt.Printf("matched %d / %d lines in %s\n",
//		col.MatchedLines, col.TotalLines, col.Duration())
package stats
