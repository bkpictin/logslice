package rotate_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/rotate"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestNew_MissingFile(t *testing.T) {
	_, err := rotate.New("/nonexistent/path/file.log", time.Second)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRotated_NoChange(t *testing.T) {
	p := writeTempFile(t, "hello\n")
	d, err := rotate.New(p, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("Rotated: %v", err)
	}
	if rotated {
		t.Error("expected no rotation for unchanged file")
	}
}

func TestRotated_Truncated(t *testing.T) {
	p := writeTempFile(t, "some content here\n")
	d, err := rotate.New(p, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Truncate the file
	if err := os.WriteFile(p, []byte(""), 0o644); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("Rotated: %v", err)
	}
	if !rotated {
		t.Error("expected rotation detected after truncation")
	}
}

func TestRotated_FileRemoved(t *testing.T) {
	p := writeTempFile(t, "data\n")
	d, err := rotate.New(p, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := os.Remove(p); err != nil {
		t.Fatalf("remove: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("Rotated on missing file should not error: %v", err)
	}
	if !rotated {
		t.Error("expected rotation detected when file removed")
	}
}

func TestReset_ClearsRotationState(t *testing.T) {
	p := writeTempFile(t, "initial content\n")
	d, err := rotate.New(p, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Truncate then reset
	if err := os.WriteFile(p, []byte("new\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := d.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	rotated, err := d.Rotated()
	if err != nil {
		t.Fatalf("Rotated: %v", err)
	}
	if rotated {
		t.Error("expected no rotation after Reset")
	}
}

func TestInterval_ReturnsConfigured(t *testing.T) {
	p := writeTempFile(t, "x\n")
	expected := 500 * time.Millisecond
	d, err := rotate.New(p, expected)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if d.Interval() != expected {
		t.Errorf("Interval() = %v, want %v", d.Interval(), expected)
	}
}
