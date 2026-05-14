// Package throttle implements a sliding-window lines-per-second rate
// limiter for log pipeline stages.
//
// Usage:
//
//	th := throttle.New(1000) // allow up to 1 000 lines/s
//	for line := range lines {
//		if !th.Allow() {
//			continue // drop line
//		}
//		process(line)
//	}
//	fmt.Println("dropped:", th.Dropped())
package throttle
