package pipeline

import (
	"testing"
)

func TestNewMaskLevel_NilWhenEmpty(t *testing.T) {
	cfg := Config{}
	m, err := newMaskLevel(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m != nil {
		t.Fatal("expected nil masker for empty config")
	}
}

func TestNewMaskLevel_InvalidLevel_ReturnsError(t *testing.T) {
	cfg := Config{MaskLevel: "VERBOSE"}
	_, err := newMaskLevel(cfg)
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNewMaskLevel_ValidLevel(t *testing.T) {
	cfg := Config{MaskLevel: "ERROR"}
	m, err := newMaskLevel(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil masker")
	}
	if !m.Enabled() {
		t.Fatal("expected masker to be enabled")
	}
}

func TestApplyMaskLevel_NilPassthrough(t *testing.T) {
	line := "ERROR something bad"
	if got := applyMaskLevel(nil, line); got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestApplyMaskLevel_MasksLine(t *testing.T) {
	cfg := Config{MaskLevel: "ERROR", MaskLevelToken: "[REDACTED]"}
	m, err := newMaskLevel(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := applyMaskLevel(m, "2024-01-01 ERROR disk full")
	want := "2024-01-01 [REDACTED] disk full"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplyMaskLevel_BelowThreshold_Unchanged(t *testing.T) {
	cfg := Config{MaskLevel: "ERROR"}
	m, err := newMaskLevel(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "INFO startup complete"
	if got := applyMaskLevel(m, line); got != line {
		t.Fatalf("expected no change, got %q", got)
	}
}
