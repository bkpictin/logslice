package maskfield

import (
	"strings"
	"testing"
)

func TestNew_NoFields_Disabled(t *testing.T) {
	m := New(nil)
	if m.Enabled() {
		t.Fatal("expected masker to be disabled with no fields")
	}
}

func TestNew_WithFields_Enabled(t *testing.T) {
	m := New([]string{"password", "token"})
	if !m.Enabled() {
		t.Fatal("expected masker to be enabled")
	}
	if len(m.Fields()) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(m.Fields()))
	}
}

func TestLine_NoFields_Passthrough(t *testing.T) {
	m := New(nil)
	line := `{"password":"secret","user":"alice"}`
	if got := m.Line(line); got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestLine_MasksJSONField(t *testing.T) {
	m := New([]string{"password"})
	line := `{"user":"alice","password":"s3cr3t"}`
	got := m.Line(line)
	if strings.Contains(got, "s3cr3t") {
		t.Fatalf("expected password to be masked, got %q", got)
	}
	if !strings.Contains(got, `"***"`) {
		t.Fatalf("expected default mask in output, got %q", got)
	}
}

func TestLine_MasksKVField(t *testing.T) {
	m := New([]string{"token"})
	line := `level=info msg="request" token=abc123xyz`
	got := m.Line(line)
	if strings.Contains(got, "abc123xyz") {
		t.Fatalf("expected token to be masked, got %q", got)
	}
	if !strings.Contains(got, "token=***") {
		t.Fatalf("expected token=*** in output, got %q", got)
	}
}

func TestLine_CustomMask(t *testing.T) {
	m := New([]string{"secret"}, WithMask("[REDACTED]"))
	line := `secret=mysecretvalue`
	got := m.Line(line)
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("expected custom mask, got %q", got)
	}
}

func TestLine_MultipleFields(t *testing.T) {
	m := New([]string{"password", "token"})
	line := `{"password":"hunter2","token":"tok_live_abc","user":"bob"}`
	got := m.Line(line)
	if strings.Contains(got, "hunter2") || strings.Contains(got, "tok_live_abc") {
		t.Fatalf("expected both fields masked, got %q", got)
	}
	if !strings.Contains(got, `"user":"bob"`) {
		t.Fatalf("expected unmasked field preserved, got %q", got)
	}
}

func TestLines_AppliedToAll(t *testing.T) {
	m := New([]string{"key"})
	input := []string{
		`key=val1 msg=hello`,
		`key=val2 msg=world`,
		`msg=plain`,
	}
	out := m.Lines(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d lines, got %d", len(input), len(out))
	}
	for _, l := range out[:2] {
		if strings.Contains(l, "val1") || strings.Contains(l, "val2") {
			t.Fatalf("expected values masked, got %q", l)
		}
	}
	if out[2] != input[2] {
		t.Fatalf("expected plain line unchanged, got %q", out[2])
	}
}

func TestReset_DisablesMasker(t *testing.T) {
	m := New([]string{"password"})
	m.Reset()
	if m.Enabled() {
		t.Fatal("expected masker disabled after Reset")
	}
	line := `password=secret`
	if got := m.Line(line); got != line {
		t.Fatalf("expected passthrough after reset, got %q", got)
	}
}
