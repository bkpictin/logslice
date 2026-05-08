package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/ratelimit"
)

func TestNew_Disabled_WhenZero(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()
	if !l.Disabled() {
		t.Fatal("expected limiter to be disabled for linesPerSecond=0")
	}
}

func TestNew_Disabled_WhenNegative(t *testing.T) {
	l := ratelimit.New(-5)
	defer l.Stop()
	if !l.Disabled() {
		t.Fatal("expected limiter to be disabled for negative linesPerSecond")
	}
}

func TestNew_Enabled(t *testing.T) {
	l := ratelimit.New(100)
	defer l.Stop()
	if l.Disabled() {
		t.Fatal("expected limiter to be enabled for linesPerSecond=100")
	}
}

func TestWait_DisabledReturnsImmediately(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 50; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("disabled limiter should not block, took %v", elapsed)
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	// Use a very slow limiter so the context cancels before any tick.
	l := ratelimit.New(1)
	defer l.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}

func TestWait_HighRate_PassesManyLines(t *testing.T) {
	l := ratelimit.New(10_000)
	defer l.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	const want = 20
	for i := 0; i < want; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}
	}
}
