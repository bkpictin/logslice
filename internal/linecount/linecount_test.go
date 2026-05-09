package linecount_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func TestNew_TotalStartsAtZero(t *testing.T) {
	c := linecount.New(10, nil)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestInc_IncreasesTotal(t *testing.T) {
	c := linecount.New(0, nil)
	c.Inc()
	c.Inc()
	c.Inc()
	if got := c.Total(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestAdd_NegativeDeltaIgnored(t *testing.T) {
	c := linecount.New(0, nil)
	c.Add(-5)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsTotal(t *testing.T) {
	c := linecount.New(0, nil)
	for i := 0; i < 20; i++ {
		c.Inc()
	}
	c.Reset()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCallback_FiredAtInterval(t *testing.T) {
	var fired []int64
	c := linecount.New(10, func(n int64) {
		fired = append(fired, n)
	})

	for i := 0; i < 25; i++ {
		c.Inc()
	}

	if len(fired) != 2 {
		t.Fatalf("expected 2 callbacks (at 10 and 20), got %d: %v", len(fired), fired)
	}
	if fired[0] != 10 || fired[1] != 20 {
		t.Fatalf("unexpected callback values: %v", fired)
	}
}

func TestCallback_NoInterval_NeverFired(t *testing.T) {
	called := false
	c := linecount.New(0, func(n int64) { called = true })
	for i := 0; i < 100; i++ {
		c.Inc()
	}
	if called {
		t.Fatal("callback should not fire when interval is 0")
	}
}

func TestAdd_BulkCrossesMultipleIntervals(t *testing.T) {
	var fired []int64
	c := linecount.New(10, func(n int64) {
		fired = append(fired, n)
	})
	c.Add(35)
	if len(fired) != 3 {
		t.Fatalf("expected 3 callbacks, got %d: %v", len(fired), fired)
	}
}

func TestCounter_ConcurrentInc(t *testing.T) {
	const goroutines = 10
	const perGoroutine = 100

	c := linecount.New(0, nil)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()

	if got := c.Total(); got != goroutines*perGoroutine {
		t.Fatalf("expected %d, got %d", goroutines*perGoroutine, got)
	}
}
