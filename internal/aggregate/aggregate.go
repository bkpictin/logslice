// Package aggregate provides log line aggregation by a key field,
// counting occurrences and optionally capping the result set.
package aggregate

import "sync"

// Entry holds the aggregated count for a single key.
type Entry struct {
	Key   string
	Count int64
}

// Aggregator accumulates counts keyed by an extracted string value.
type Aggregator struct {
	mu      sync.Mutex
	counts  map[string]int64
	order   []string // insertion order
	maxKeys int
}

// New creates an Aggregator. maxKeys limits the number of distinct keys
// tracked; zero or negative means unlimited.
func New(maxKeys int) *Aggregator {
	return &Aggregator{
		counts:  make(map[string]int64),
		maxKeys: maxKeys,
	}
}

// Add records one occurrence of key. If maxKeys is set and already reached,
// the key is silently dropped unless it already exists.
func (a *Aggregator) Add(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.counts[key]; !exists {
		if a.maxKeys > 0 && len(a.order) >= a.maxKeys {
			return
		}
		a.order = append(a.order, key)
	}
	a.counts[key]++
}

// Results returns aggregated entries in insertion order.
func (a *Aggregator) Results() []Entry {
	a.mu.Lock()
	defer a.mu.Unlock()

	out := make([]Entry, 0, len(a.order))
	for _, k := range a.order {
		out = append(out, Entry{Key: k, Count: a.counts[k]})
	}
	return out
}

// Reset clears all accumulated state.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.counts = make(map[string]int64)
	a.order = a.order[:0]
}

// Len returns the number of distinct keys currently tracked.
func (a *Aggregator) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.counts)
}

// TopN returns up to n entries sorted by count descending. Entries with equal
// counts retain their original insertion order. If n is zero or negative, all
// entries are returned.
func (a *Aggregator) TopN(n int) []Entry {
	entries := a.Results()

	// Stable sort by count descending, preserving insertion order for ties.
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Count > entries[j-1].Count; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}

	if n > 0 && n < len(entries) {
		return entries[:n]
	}
	return entries
}
