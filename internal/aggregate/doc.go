// Package aggregate provides a thread-safe log-line aggregator that groups
// lines by an extracted key (e.g. a log level, hostname, or field value) and
// counts how many times each key appears.
//
// Typical usage:
//
//	agg := aggregate.New(100) // track at most 100 distinct keys
//	for _, line := range lines {
//		agg.Add(extractKey(line))
//	}
//	for _, e := range agg.Results() {
//		fmt.Printf("%s\t%d\n", e.Key, e.Count)
//	}
package aggregate
