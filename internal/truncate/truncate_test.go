package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/truncate"
)

func TestLine_NoTruncation(t *testing.T) {
	tr := truncate.New(0, "...")
	s := strings.Repeat("a", 200)
	if got := tr.Line(s); got != s {
		t.Fatalf("expected unchanged line, got len %d", len(got))
	}
}

func TestLine_ShortEnough(t *testing.T) {
	tr := truncate.NewDefault(100)
	s := "short line"
	if got := tr.Line(s); got != s {
		t.Fatalf("expected %q, got %q", s, got)
	}
}

func TestLine_TruncatesWithSuffix(t *testing.T) {
	tr := truncate.NewDefault(10)
	s := "0123456789EXTRA"
	got := tr.Line(s)
	if len(got) != 10 {
		t.Fatalf("expected len 10, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestLine_TruncatesNoSuffix(t *testing.T) {
	tr := truncate.New(5, "")
	got := tr.Line("abcdefgh")
	if got != "abcde" {
		t.Fatalf("expected %q, got %q", "abcde", got)
	}
}

func TestLines_AppliedToAll(t *testing.T) {
	tr := truncate.NewDefault(6)
	input := []string{"short", "toolongline", "ok"}
	out := tr.Lines(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 results")
	}
	if out[1] != "too..." {
		t.Fatalf("expected 'too...', got %q", out[1])
	}
	if out[0] != "short" {
		t.Fatalf("short line should be unchanged")
	}
}

func TestEnabled(t *testing.T) {
	if truncate.New(0, "").Enabled() {
		t.Fatal("expected Enabled() == false when maxLen==0")
	}
	if !truncate.New(80, "").Enabled() {
		t.Fatal("expected Enabled() == true when maxLen>0")
	}
}

func TestTrimFields(t *testing.T) {
	s := "hello worldlongtoken end"
	got := truncate.TrimFields(s, 8)
	if strings.Contains(got, "worldlongtoken") {
		t.Fatalf("long field was not trimmed: %q", got)
	}
}

func TestTrimFields_Disabled(t *testing.T) {
	s := "hello worldlongtoken"
	if got := truncate.TrimFields(s, 0); got != s {
		t.Fatalf("expected unchanged, got %q", got)
	}
}
