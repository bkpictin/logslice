// Package merge provides ordered merging of log lines from multiple
// sorted sources, emitting lines in chronological order based on timestamp.
package merge

import (
	"container/heap"
	"time"
)

// Line represents a single log line with an associated timestamp and
// source identifier for tracking origin after merging.
type Line struct {
	Text      string
	Timestamp time.Time
	Source    string
	LineNum   int
}

// entry is an internal heap element.
type entry struct {
	line  Line
	index int // source index
}

// minHeap implements heap.Interface for chronological ordering.
type minHeap []entry

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].line.Timestamp.Before(h[j].line.Timestamp) }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(entry)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges multiple ordered line channels into a single ordered stream.
type Merger struct {
	sources []<-chan Line
}

// New creates a Merger that reads from the provided source channels.
// Each channel must emit lines in non-decreasing timestamp order.
func New(sources []<-chan Line) *Merger {
	return &Merger{sources: sources}
}

// Merge reads from all sources and emits lines in timestamp order on the
// returned channel. The channel is closed when all sources are exhausted.
func (m *Merger) Merge() <-chan Line {
	out := make(chan Line, 64)
	go func() {
		defer close(out)
		h := &minHeap{}
		heap.Init(h)

		chans := make([]<-chan Line, len(m.sources))
		copy(chans, m.sources)

		// Seed the heap with the first line from each source.
		for i, ch := range chans {
			if l, ok := <-ch; ok {
				heap.Push(h, entry{line: l, index: i})
			}
		}

		for h.Len() > 0 {
			e := heap.Pop(h).(entry)
			out <- e.line
			if l, ok := <-chans[e.index]; ok {
				heap.Push(h, entry{line: l, index: e.index})
			}
		}
	}()
	return out
}
