// Package highlight provides ANSI terminal colour support for logslice output.
//
// It is intentionally kept separate from the output formatter so that colour
// rendering can be toggled at runtime (e.g. when stdout is not a TTY or when
// the user passes --no-color) without affecting structured output formats such
// as JSON or CSV.
//
// Usage:
//
//	h := highlight.New(isTTY, compiledPattern)
//	colouredLine := h.Line(rawLine, "ERROR")
package highlight
