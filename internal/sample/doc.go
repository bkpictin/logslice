// Package sample implements deterministic 1-in-N log line sampling for
// logslice pipelines.
//
// When processing very high-volume log files it is often useful to inspect
// only a representative fraction of the lines rather than every entry.
// The Sampler type provides a lightweight, goroutine-safe counter that
// approves every Nth line and silently drops the rest, keeping the output
// manageable without losing the overall shape of the data.
//
// Usage:
//
//	s := sample.New(10)   // keep 1 in every 10 lines
//	for _, line := range lines {
//	    if s.Keep() {
//	        emit(line)
//	    }
//	}
//	fmt.Println("dropped:", s.Skipped())
package sample
