package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/reader"
)

func makeLinesChan(texts []string) <-chan reader.Line {
	ch := make(chan reader.Line, len(texts))
	for i, t := range texts {
		ch <- reader.Line{Text: t, Number: i + 1}
	}
	close(ch)
	return ch
}

func collectLines(ch <-chan reader.Line, timeout time.Duration) []string {
	var out []string
	timer := time.After(timeout)
	for {
		select {
		case l, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, l.Text)
		case <-timer:
			return out
		}
	}
}

func TestNewMultilineCollector_NilWhenNoPattern(t *testing.T) {
	cfg := Config{}
	c, err := newMultilineCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c != nil {
		t.Fatal("expected nil collector when MultilineStart is empty")
	}
}

func TestNewMultilineCollector_InvalidPattern(t *testing.T) {
	cfg := Config{MultilineStart: "[invalid("}
	_, err := newMultilineCollector(cfg)
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestApplyMultiline_NilPassthrough(t *testing.T) {
	ctx := context.Background()
	input := makeLinesChan([]string{"line1", "line2"})
	out := applyMultiline(ctx, input, nil)
	got := collectLines(out, time.Second)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestApplyMultiline_FoldsLines(t *testing.T) {
	cfg := Config{MultilineStart: `^\d{4}-`}
	collector, err := newMultilineCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := makeLinesChan([]string{
		"2024-01-01 start of entry",
		"  continuation line",
		"2024-01-02 second entry",
	})
	ctx := context.Background()
	out := applyMultiline(ctx, input, collector)
	got := collectLines(out, time.Second)
	if len(got) != 2 {
		t.Fatalf("expected 2 folded entries, got %d: %v", len(got), got)
	}
}

func TestApplyMultiline_ContextCancel(t *testing.T) {
	cfg := Config{MultilineStart: `^START`}
	collector, err := newMultilineCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	blocking := make(chan reader.Line) // never closed
	out := applyMultiline(ctx, blocking, collector)
	cancel()
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed after cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close after cancel")
	}
}
