package burst_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/burst"
)

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNew_DisabledWhenZero(t *testing.T) {
	d := burst.New(0, time.Minute)
	for i := 0; i < 1000; i++ {
		if d.Observe(base) {
			t.Fatal("expected disabled detector to never trigger")
		}
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	d := burst.New(-1, time.Minute)
	if d.Observe(base) {
		t.Fatal("expected disabled detector to never trigger")
	}
}

func TestObserve_BelowThreshold(t *testing.T) {
	d := burst.New(5, time.Minute)
	for i := 0; i < 5; i++ {
		if d.Observe(base) {
			t.Fatalf("triggered early at event %d", i+1)
		}
	}
}

func TestObserve_ExceedsThreshold(t *testing.T) {
	d := burst.New(3, time.Minute)
	for i := 0; i < 3; i++ {
		d.Observe(base)
	}
	if !d.Observe(base) {
		t.Fatal("expected burst to be detected on 4th event")
	}
}

func TestObserve_SlidingWindowEvicts(t *testing.T) {
	d := burst.New(3, time.Second)
	// Fill window at t=0
	for i := 0; i < 3; i++ {
		d.Observe(base)
	}
	// Advance past the window; old events should be evicted
	now := base.Add(2 * time.Second)
	if d.Observe(now) {
		t.Fatal("expected no burst after window slides past old events")
	}
}

func TestTriggeredCount(t *testing.T) {
	d := burst.New(2, time.Minute)
	d.Observe(base)
	d.Observe(base)
	d.Observe(base) // 3rd — triggers
	d.Observe(base) // 4th — triggers
	if got := d.TriggeredCount(); got != 2 {
		t.Fatalf("expected 2 triggers, got %d", got)
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := burst.New(2, time.Minute)
	d.Observe(base)
	d.Observe(base)
	d.Observe(base) // triggers
	d.Reset()
	if d.TriggeredCount() != 0 {
		t.Fatal("expected triggered count to be 0 after reset")
	}
	// Should not trigger immediately after reset
	if d.Observe(base) {
		t.Fatal("expected no burst immediately after reset")
	}
}
