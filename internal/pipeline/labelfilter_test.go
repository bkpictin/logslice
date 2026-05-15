package pipeline

import (
	"testing"
)

func TestNewLabelFilter_NilWhenNoLabels(t *testing.T) {
	cfg := Config{}
	f, err := newLabelFilter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Error("expected nil filter when no labels configured")
	}
}

func TestNewLabelFilter_InvalidPair_ReturnsError(t *testing.T) {
	cfg := Config{Labels: []string{"noeq"}}
	_, err := newLabelFilter(cfg)
	if err == nil {
		t.Fatal("expected error for invalid label pair")
	}
}

func TestNewLabelFilter_ValidPairs(t *testing.T) {
	cfg := Config{Labels: []string{"env=prod", "region=us"}}
	f, err := newLabelFilter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if !f.Enabled() {
		t.Error("expected filter to be enabled")
	}
}

func TestApplyLabelFilter_NilPassthrough(t *testing.T) {
	if !applyLabelFilter(nil, "any line") {
		t.Error("nil filter should pass every line")
	}
}

func TestApplyLabelFilter_DropsNonMatching(t *testing.T) {
	cfg := Config{Labels: []string{"env=prod"}}
	f, _ := newLabelFilter(cfg)
	if applyLabelFilter(f, "env=staging msg=hi") {
		t.Error("expected line to be dropped")
	}
}

func TestApplyLabelFilter_PassesMatching(t *testing.T) {
	cfg := Config{Labels: []string{"env=prod"}}
	f, _ := newLabelFilter(cfg)
	if !applyLabelFilter(f, "env=prod msg=hi") {
		t.Error("expected line to pass")
	}
}
