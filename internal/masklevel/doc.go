// Package masklevel provides a line transformer that replaces log-level tokens
// with a configurable mask string.
//
// It is useful when forwarding log lines to downstream systems that require a
// normalised severity field, or when severity information must be redacted from
// output for compliance reasons.
//
// Usage:
//
//	m, err := masklevel.New("WARN", masklevel.WithMask("***"))
//	if err != nil { /* handle */ }
//	out := m.Line("2024-01-01 ERROR something went wrong")
//	// out == "2024-01-01 *** something went wrong"
package masklevel
