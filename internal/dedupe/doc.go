// Package dedupe implements a sliding-window deduplication filter for
// log lines. It is designed to be embedded in a pipeline stage so that
// repeated log entries (e.g. a tight error loop) do not flood output.
//
// Usage:
//
//	d := dedupe.New(10)          // remember last 10 distinct lines
//	for _, line := range lines {
//		if !d.IsDuplicate(line) {
//			fmt.Println(line)
//		}
//	}
//	fmt.Printf("%d duplicates suppressed\n", d.Skipped)
package dedupe
