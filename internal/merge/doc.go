// Package merge implements a k-way merge of sorted log line streams.
//
// It accepts multiple channels of Line values, each expected to emit lines
// in non-decreasing timestamp order, and produces a single output channel
// that yields lines in globally sorted order using a min-heap.
//
// Typical usage:
//
//	ch1 := makeSource("app.log")
//	ch2 := makeSource("worker.log")
//	m := merge.New([]<-chan merge.Line{ch1, ch2})
//	for line := range m.Merge() {
//		fmt.Println(line.Source, line.Text)
//	}
package merge
