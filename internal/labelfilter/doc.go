// Package labelfilter provides a key=value label matcher for log lines.
//
// A Filter is constructed from a list of "key=value" pairs and reports
// whether a given log line contains all of them.  This is useful for
// filtering structured log output (e.g. logfmt or inline labels) without
// requiring full JSON parsing.
//
// Example:
//
//	f, _ := labelfilter.New([]string{"env=prod", "service=api"})
//	if f.Match(line) {
//		// line contains both env=prod and service=api
//	}
package labelfilter
