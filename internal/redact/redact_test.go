package redact

import (
	"testing"
)

func TestNew_NoPatterns_Disabled(t *testing.T) {
	r, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Enabled() {
		t.Error("expected Enabled() == false with no patterns")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]string{`[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestLine_NoPatterns_Passthrough(t *testing.T) {
	r, _ := New(nil)
	line := "user logged in with password=secret"
	if got := r.Line(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestLine_RedactsMatch(t *testing.T) {
	r, err := New([]string{`password=\S+`})
	if err != nil {
		t.Fatal(err)
	}
	input := "auth failed password=s3cr3t for user=alice"
	got := r.Line(input)
	want := "auth failed [REDACTED] for user=alice"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLine_CustomPlaceholder(t *testing.T) {
	r, err := New([]string{`token=[A-Za-z0-9]+`}, WithPlaceholder("***"))
	if err != nil {
		t.Fatal(err)
	}
	got := r.Line("request token=abc123 received")
	want := "request *** received"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLine_MultiplePatterns(t *testing.T) {
	r, err := New([]string{`password=\S+`, `secret=\S+`})
	if err != nil {
		t.Fatal(err)
	}
	got := r.Line("password=abc secret=xyz other=ok")
	want := "[REDACTED] [REDACTED] other=ok"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLines_AppliedToAll(t *testing.T) {
	r, err := New([]string{`key=\S+`})
	if err != nil {
		t.Fatal(err)
	}
	input := []string{"key=abc", "no match", "key=xyz end"}
	out := r.Lines(input)
	expected := []string{"[REDACTED]", "no match", "[REDACTED] end"}
	for i, want := range expected {
		if out[i] != want {
			t.Errorf("line %d: got %q, want %q", i, out[i], want)
		}
	}
}

func TestLines_Disabled_ReturnsOriginal(t *testing.T) {
	r, _ := New(nil)
	input := []string{"line one", "line two"}
	out := r.Lines(input)
	if &out[0] == &input[0] {
		// same slice returned — fine
	}
	for i := range input {
		if out[i] != input[i] {
			t.Errorf("line %d changed unexpectedly", i)
		}
	}
}
