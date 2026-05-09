// Package burst implements a sliding-window burst detector for logslice.
//
// A Detector counts how many log lines are observed within a configurable
// time window. When the count exceeds the configured threshold the Observe
// call returns true, signalling that a burst of activity has been detected.
//
// Typical usage:
//
//	det := burst.New(100, time.Minute)
//	if det.Observe(time.Now()) {
//		log.Println("burst detected")
//	}
package burst
