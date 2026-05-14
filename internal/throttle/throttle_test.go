package throttle

import (
	"testing"
	"time"
)

func TestNew_DisabledWhenZero(t *testing.T) {
	th := New(0)
	if th.Enabled() {
		t.Fatal("expected disabled for maxLPS=0")
	}
	for i := 0; i < 1000; i++ {
		if !th.Allow() {
			t.Fatal("disabled throttler should always allow")
		}
	}
}

func TestNew_DisabledWhenNegative(t *testing.T) {
	th := New(-5)
	if th.Enabled() {
		t.Fatal("expected disabled for negative maxLPS")
	}
}

func TestAllow_UnderLimit(t *testing.T) {
	th := New(100)
	for i := 0; i < 100; i++ {
		if !th.Allow() {
			t.Fatalf("line %d should be allowed (under limit)", i)
		}
	}
	if th.Dropped() != 0 {
		t.Fatalf("expected 0 dropped, got %d", th.Dropped())
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	const limit = 5
	th := New(limit)
	allowed := 0
	for i := 0; i < limit+10; i++ {
		if th.Allow() {
			allowed++
		}
	}
	if allowed != limit {
		t.Fatalf("expected %d allowed, got %d", limit, allowed)
	}
	if th.Dropped() != 10 {
		t.Fatalf("expected 10 dropped, got %d", th.Dropped())
	}
}

func TestAllow_SlidingWindowRefills(t *testing.T) {
	th := New(3)
	for i := 0; i < 3; i++ {
		th.Allow()
	}
	// bucket should be full; next call drops
	if th.Allow() {
		t.Fatal("expected drop when bucket full")
	}
	// Manually age the bucket entries by manipulating internal state
	th.mu.Lock()
	old := time.Now().Add(-2 * time.Second)
	for i := range th.buckets {
		th.buckets[i] = old
	}
	th.mu.Unlock()

	// Window has expired; should allow again
	if !th.Allow() {
		t.Fatal("expected allow after window expiry")
	}
}

func TestReset_ClearsState(t *testing.T) {
	th := New(2)
	th.Allow()
	th.Allow()
	th.Allow() // dropped
	if th.Dropped() == 0 {
		t.Fatal("expected at least one drop before reset")
	}
	th.Reset()
	if th.Dropped() != 0 {
		t.Fatalf("expected 0 dropped after reset, got %d", th.Dropped())
	}
	// should allow again after reset
	if !th.Allow() {
		t.Fatal("expected allow after reset")
	}
}
