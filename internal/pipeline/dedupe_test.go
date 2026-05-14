package pipeline

import (
	"testing"

	"github.com/user/logslice/internal/filter"
)

func TestNewDeduplicator_Disabled(t *testing.T) {
	cfg := Config{DedupeWindow: 0}
	d, err := newDeduplicator(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != nil {
		t.Fatal("expected nil deduplicator when window is 0")
	}
}

func TestNewDeduplicator_Enabled(t *testing.T) {
	cfg := Config{DedupeWindow: 5}
	d, err := newDeduplicator(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil deduplicator when window > 0")
	}
}

func TestApplyDedupe_NilPassthrough(t *testing.T) {
	entry := filter.Entry{Line: "hello world"}
	got, keep := applyDedupe(nil, entry)
	if !keep {
		t.Fatal("nil deduplicator should keep all lines")
	}
	if got.Line != entry.Line {
		t.Fatalf("entry mutated unexpectedly: %q", got.Line)
	}
}

func TestApplyDedupe_DropsDuplicate(t *testing.T) {
	cfg := Config{DedupeWindow: 10}
	d, _ := newDeduplicator(cfg)

	entry := filter.Entry{Line: "duplicate line"}

	_, keep1 := applyDedupe(d, entry)
	if !keep1 {
		t.Fatal("first occurrence should be kept")
	}

	_, keep2 := applyDedupe(d, entry)
	if keep2 {
		t.Fatal("second occurrence should be dropped as duplicate")
	}
}

func TestApplyDedupe_KeepsDistinctLines(t *testing.T) {
	cfg := Config{DedupeWindow: 10}
	d, _ := newDeduplicator(cfg)

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		_, keep := applyDedupe(d, filter.Entry{Line: l})
		if !keep {
			t.Fatalf("distinct line %q should be kept", l)
		}
	}
}

func TestApplyDedupe_NegativeWindowDisabled(t *testing.T) {
	cfg := Config{DedupeWindow: -1}
	d, err := newDeduplicator(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != nil {
		t.Fatal("expected nil deduplicator when window is negative")
	}
}
