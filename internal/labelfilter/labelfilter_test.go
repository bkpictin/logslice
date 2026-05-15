package labelfilter

import (
	"testing"
)

func TestNew_NoPairs_Disabled(t *testing.T) {
	f, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Enabled() {
		t.Error("expected filter to be disabled")
	}
}

func TestNew_InvalidPair_ReturnsError(t *testing.T) {
	_, err := New([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for invalid pair")
	}
}

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestMatch_Disabled_AlwaysTrue(t *testing.T) {
	f, _ := New([]string{})
	if !f.Match("anything here") {
		t.Error("disabled filter should always match")
	}
}

func TestMatch_SingleLabel_Found(t *testing.T) {
	f, _ := New([]string{"env=prod"})
	if !f.Match("ts=2024-01-01 env=prod msg=hello") {
		t.Error("expected match")
	}
}

func TestMatch_SingleLabel_NotFound(t *testing.T) {
	f, _ := New([]string{"env=prod"})
	if f.Match("ts=2024-01-01 env=staging msg=hello") {
		t.Error("expected no match")
	}
}

func TestMatch_MultipleLabels_AllPresent(t *testing.T) {
	f, _ := New([]string{"env=prod", "service=api"})
	if !f.Match("env=prod service=api level=info msg=ok") {
		t.Error("expected match")
	}
}

func TestMatch_MultipleLabels_OneMissing(t *testing.T) {
	f, _ := New([]string{"env=prod", "service=api"})
	if f.Match("env=prod service=worker level=info msg=ok") {
		t.Error("expected no match")
	}
}

func TestCounters_MatchedAndRejected(t *testing.T) {
	f, _ := New([]string{"env=prod"})
	f.Match("env=prod a=1")
	f.Match("env=prod b=2")
	f.Match("env=staging c=3")
	if f.Matched() != 2 {
		t.Errorf("matched: want 2, got %d", f.Matched())
	}
	if f.Rejected() != 1 {
		t.Errorf("rejected: want 1, got %d", f.Rejected())
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	f, _ := New([]string{"env=prod"})
	f.Match("env=prod")
	f.Reset()
	if f.Matched() != 0 || f.Rejected() != 0 {
		t.Error("expected counters to be zero after reset")
	}
}

func TestParseError_Message(t *testing.T) {
	err := &ParseError{Raw: "bad"}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}
