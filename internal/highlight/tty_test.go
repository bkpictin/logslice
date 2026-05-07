package highlight_test

import (
	"os"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

// TestAutoEnabled_NoColor verifies that setting NO_COLOR disables highlighting
// regardless of whether stdout is a terminal.
func TestAutoEnabled_NoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if highlight.AutoEnabled() {
		t.Error("expected AutoEnabled to return false when NO_COLOR is set")
	}
}

// TestIsTTY_RegularFile ensures that a plain file is not reported as a TTY.
func TestIsTTY_RegularFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tty-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if highlight.IsTTY(f) {
		t.Error("expected regular file to not be a TTY")
	}
}
