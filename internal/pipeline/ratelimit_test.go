package pipeline

import (
	"context"
	"testing"
	"time"
)

func TestNewRateLimiter_NilWhenZero(t *testing.T) {
	cfg := Config{LinesPerSecond: 0}
	rl, err := newRateLimiter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl != nil {
		t.Fatalf("expected nil limiter for zero rate, got non-nil")
	}
}

func TestNewRateLimiter_NilWhenNegative(t *testing.T) {
	cfg := Config{LinesPerSecond: -5}
	rl, err := newRateLimiter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl != nil {
		t.Fatalf("expected nil limiter for negative rate, got non-nil")
	}
}

func TestNewRateLimiter_NonNilWhenPositive(t *testing.T) {
	cfg := Config{LinesPerSecond: 100}
	rl, err := newRateLimiter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl == nil {
		t.Fatal("expected non-nil limiter for positive rate")
	}
}

func TestApplyRateLimit_NilPassthrough(t *testing.T) {
	line, err := applyRateLimit(context.Background(), nil, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line != "hello" {
		t.Fatalf("expected \"hello\", got %q", line)
	}
}

func TestApplyRateLimit_ContextCancelled(t *testing.T) {
	cfg := Config{LinesPerSecond: 1}
	rl, err := newRateLimiter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Exhaust the first token so the next call must wait.
	_, _ = applyRateLimit(ctx, rl, "first")

	// The context should expire before the limiter grants the next token.
	_, gotErr := applyRateLimit(ctx, rl, "second")
	if gotErr == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestApplyRateLimit_HighRate_PassesThrough(t *testing.T) {
	cfg := Config{LinesPerSecond: 100000}
	rl, err := newRateLimiter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 10; i++ {
		line, err := applyRateLimit(context.Background(), rl, "line")
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		if line != "line" {
			t.Fatalf("iteration %d: expected \"line\", got %q", i, line)
		}
	}
}
