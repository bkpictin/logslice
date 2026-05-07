package highlight

import (
	"os"

	"golang.org/x/term"
)

// IsTTY reports whether f is connected to a terminal.
// It is used by callers to decide whether to enable colour output.
func IsTTY(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// AutoEnabled returns true when stdout is a TTY and the NO_COLOR environment
// variable is not set (see https://no-color.org).
func AutoEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return IsTTY(os.Stdout)
}
