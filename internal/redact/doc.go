// Package redact provides pattern-based log line redaction.
//
// A Redactor holds a set of compiled regular expressions and replaces
// any matching substrings with a configurable placeholder string
// (default "[REDACTED]"). This is useful for stripping passwords,
// API keys, tokens, or other sensitive values before writing log
// output to disk or forwarding it downstream.
//
// Usage:
//
//	r, err := redact.New([]string{`password=\S+`, `token=[A-Za-z0-9]+`})
//	if err != nil { ... }
//	clean := r.Line(rawLine)
package redact
