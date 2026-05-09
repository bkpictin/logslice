package aggregate_test

import (
	"sync"
	"testing"

	"github.com/logslice/logslice/internal/aggregate"
)

func TestNew_InitialState(t *testing.T) {
	a := aggregate.New(0)
	if a.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", a.Len())
	}
	if results := a.Results(); len(results) != 0 {
		t.Fatalf("expected empty results, got %v", results)
	}
}

func TestAdd_CountsCorrectly(t *testing.T) {
	a := aggregate.New(0)
	a.Add("error")
	a.Add("warn")
	a.Add("error")
	a.Add("error")

	results := a.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(results))
	}
	if results[0].Key != "error" || results[0].Count != 3 {
		t.Errorf("unexpected first entry: %+v", results[0])
	}
	if results[1].Key != "warn" || results[1].Count != 1 {
		t.Errorf("unexpected second entry: %+v", results[1])
	}
}

func TestAdd_InsertionOrder(t *testing.T) {
	a := aggregate.New(0)
	keys := []string{"c", "a", "b", "a", "c"}
	for _, k := range keys {
		a.Add(k)
	}
	results := a.Results()
	expected := []string{"c", "a", "b"}
	for i, e := range results {
		if e.Key != expected[i] {
			t.Errorf("position %d: want %q got %q", i, expected[i], e.Key)
		}
	}
}

func TestAdd_MaxKeys_DropsNewKeys(t *testing.T) {
	a := aggregate.New(2)
	a.Add("alpha")
	a.Add("beta")
	a.Add("gamma") // should be dropped
	a.Add("alpha") // existing key — must still be counted

	if a.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", a.Len())
	}
	for _, e := range a.Results() {
		if e.Key == "gamma" {
			t.Error("gamma should have been dropped")
		}
		if e.Key == "alpha" && e.Count != 2 {
			t.Errorf("alpha count: want 2, got %d", e.Count)
		}
	}
}

func TestReset_ClearsState(t *testing.T) {
	a := aggregate.New(0)
	a.Add("x")
	a.Add("y")
	a.Reset()
	if a.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", a.Len())
	}
}

func TestAdd_ConcurrentSafe(t *testing.T) {
	a := aggregate.New(0)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			a.Add("key")
		}(i)
	}
	wg.Wait()
	results := a.Results()
	if len(results) != 1 || results[0].Count != 100 {
		t.Errorf("expected 1 key with count 100, got %+v", results)
	}
}
