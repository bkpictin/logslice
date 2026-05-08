package dedupe_test

import (
	"fmt"
	"testing"

	"github.com/user/logslice/internal/dedupe"
)

func TestNew_WindowOne(t *testing.T) {
	d := dedupe.New(1)
	if d.IsDuplicate("hello") {
		t.Fatal("first occurrence should not be a duplicate")
	}
	if !d.IsDuplicate("hello") {
		t.Fatal("immediate repeat should be a duplicate")
	}
	if d.IsDuplicate("world") {
		t.Fatal("different line should not be a duplicate")
	}
}

func TestNew_WindowEvictsOldEntries(t *testing.T) {
	d := dedupe.New(3)
	lines := []string{"a", "b", "c", "d"}
	for _, l := range lines {
		if d.IsDuplicate(l) {
			t.Fatalf("unexpected duplicate for %q", l)
		}
	}
	// "a" should have been evicted by the fourth insertion
	if d.IsDuplicate("a") {
		t.Fatal("'a' should have been evicted from the window")
	}
}

func TestSkippedCount(t *testing.T) {
	d := dedupe.New(5)
	for i := 0; i < 4; i++ {
		d.IsDuplicate("repeat")
	}
	if d.Skipped != 3 {
		t.Fatalf("expected 3 skipped, got %d", d.Skipped)
	}
}

func TestSeen_CountsUnique(t *testing.T) {
	d := dedupe.New(10)
	for i := 0; i < 5; i++ {
		d.IsDuplicate(fmt.Sprintf("line-%d", i))
	}
	d.IsDuplicate("line-0") // duplicate
	if d.Seen() != 5 {
		t.Fatalf("expected 5 unique lines, got %d", d.Seen())
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := dedupe.New(4)
	d.IsDuplicate("foo")
	d.IsDuplicate("foo")
	d.Reset()
	if d.Skipped != 0 {
		t.Fatalf("expected Skipped=0 after Reset, got %d", d.Skipped)
	}
	if d.Seen() != 0 {
		t.Fatalf("expected Seen=0 after Reset, got %d", d.Seen())
	}
	if d.IsDuplicate("foo") {
		t.Fatal("'foo' should not be a duplicate after Reset")
	}
}

func TestNew_InvalidWindowDefaultsToOne(t *testing.T) {
	d := dedupe.New(0)
	d.IsDuplicate("x")
	if !d.IsDuplicate("x") {
		t.Fatal("window=0 should default to 1; repeat should be duplicate")
	}
}
