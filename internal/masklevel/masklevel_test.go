package masklevel_test

import (
	"testing"

	"github.com/logslice/logslice/internal/masklevel"
)

func TestNew_EmptyLevel_Disabled(t *testing.T) {
	m, err := masklevel.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Enabled() {
		t.Fatal("expected masker to be disabled for empty level")
	}
}

func TestNew_UnknownLevel_ReturnsError(t *testing.T) {
	_, err := masklevel.New("VERBOSE")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestLine_Disabled_Passthrough(t *testing.T) {
	m, _ := masklevel.New("")
	line := "INFO something happened"
	if got := m.Line(line); got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestLine_MasksMatchingLevel(t *testing.T) {
	m, err := masklevel.New("ERROR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Line("2024-01-01 ERROR disk full")
	if got != "2024-01-01 [LEVEL] disk full" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestLine_DoesNotMaskBelowMinLevel(t *testing.T) {
	m, err := masklevel.New("ERROR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "2024-01-01 INFO startup complete"
	if got := m.Line(line); got != line {
		t.Fatalf("expected no change for INFO line, got %q", got)
	}
}

func TestLine_CustomMask(t *testing.T) {
	m, err := masklevel.New("WARN", masklevel.WithMask("***"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Line("WARN threshold exceeded")
	if got != "*** threshold exceeded" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestLine_NoLevelToken_Unchanged(t *testing.T) {
	m, err := masklevel.New("WARN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "plain log line without level"
	if got := m.Line(line); got != line {
		t.Fatalf("expected no change, got %q", got)
	}
}
