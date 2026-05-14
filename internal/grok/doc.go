// Package grok provides Grok-style log pattern matching.
//
// Grok patterns are named regular-expression templates that can be composed
// using %{PATTERN_NAME} or %{PATTERN_NAME:field_name} syntax.  Matching a
// line returns a map of named field captures, making it easy to extract
// structured data from unstructured log lines without writing raw regexps.
//
// A small set of built-in patterns (IP, LOGLEVEL, TIMESTAMP, …) is included.
// Additional patterns can be registered via WithPattern / WithPatterns options.
//
// Example:
//
//	p, err := grok.New(`%{IP:client} %{WORD:method} %{NOTSPACE:path}`)
//	if err != nil { … }
//	fields := p.Match("192.168.1.1 GET /api/v1/users")
//	// fields["client"] == "192.168.1.1"
package grok
