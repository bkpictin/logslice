package lineformat

import (
	"strings"
	"testing"
)

func TestLine_TrimsCRLF(t *testing.T) {
	f := NewDefault()
	got := f.Line("hello world\r\n")
	if got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestLine_TrimsTrailingSpaces(t *testing.T) {
	f := NewDefault()
	got := f.Line("hello   ")
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestLine_StripsANSI(t *testing.T) {
	f := NewDefault()
	coloured := "\x1b[32mINFO\x1b[0m starting up"
	got := f.Line(coloured)
	if strings.Contains(got, "\x1b") {
		t.Fatalf("ANSI codes not stripped: %q", got)
	}
	if got != "INFO starting up" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestLine_NoStripANSI_WhenDisabled(t *testing.T) {
	f := New(Options{StripANSI: false})
	coloured := "\x1b[32mINFO\x1b[0m msg"
	got := f.Line(coloured)
	if !strings.Contains(got, "\x1b") {
		t.Fatalf("expected ANSI codes to be preserved, got %q", got)
	}
}

func TestLine_MaxBytes_Caps(t *testing.T) {
	f := New(Options{MaxBytes: 5})
	got := f.Line("hello world")
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestLine_MaxBytes_Zero_NoLimit(t *testing.T) {
	f := New(Options{MaxBytes: 0})
	long := strings.Repeat("x", 200)
	got := f.Line(long)
	if len(got) != 200 {
		t.Fatalf("expected 200 bytes, got %d", len(got))
	}
}

func TestLines_AppliedToAll(t *testing.T) {
	f := NewDefault()
	input := []string{"\x1b[31mERROR\x1b[0m\r\n", "ok\r", "fine"}
	out := f.Lines(input)
	expected := []string{"ERROR", "ok", "fine"}
	for i, e := range expected {
		if out[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, out[i])
		}
	}
}

func TestEnabled_WhenStripANSI(t *testing.T) {
	f := New(Options{StripANSI: true})
	if !f.Enabled() {
		t.Fatal("expected Enabled() == true")
	}
}

func TestEnabled_WhenMaxBytes(t *testing.T) {
	f := New(Options{MaxBytes: 100})
	if !f.Enabled() {
		t.Fatal("expected Enabled() == true")
	}
}

func TestEnabled_WhenNeitherSet(t *testing.T) {
	f := New(Options{})
	if f.Enabled() {
		t.Fatal("expected Enabled() == false")
	}
}
