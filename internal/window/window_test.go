package window

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func ts(sec int) time.Time { return epoch.Add(time.Duration(sec) * time.Second) }

func TestFeed_NoMatch_NothingEmitted(t *testing.T) {
	w := New(2*time.Second, 2*time.Second)
	for i := 0; i < 5; i++ {
		out := w.Feed("line", ts(i), false)
		if len(out) != 0 {
			t.Fatalf("expected no output before match, got %v", out)
		}
	}
}

func TestFeed_MatchEmitsSurroundingLines(t *testing.T) {
	w := New(2*time.Second, 2*time.Second)
	lines := []string{"a", "b", "c", "d", "e"}
	var all []string
	for i, l := range lines {
		isAnchor := l == "c"
		out := w.Feed(l, ts(i), isAnchor)
		all = append(all, out...)
	}
	// drain remaining after window closes
	out := w.Feed("f", ts(10), false)
	all = append(all, out...)

	expected := []string{"a", "b", "c", "d", "e"}
	if len(all) != len(expected) {
		t.Fatalf("expected %v got %v", expected, all)
	}
	for i, v := range expected {
		if all[i] != v {
			t.Errorf("pos %d: want %q got %q", i, v, all[i])
		}
	}
}

func TestFeed_BeforeWindowPrunesOldEntries(t *testing.T) {
	w := New(1*time.Second, 0)
	w.Feed("old", ts(0), false)
	w.Feed("recent", ts(5), false)
	out := w.Feed("anchor", ts(6), true)
	// "old" at t=0 is outside 1s before t=6 (cutoff t=5)
	for _, l := range out {
		if l == "old" {
			t.Error("old line should have been pruned")
		}
	}
}

func TestFeed_AfterWindowLimitsEmission(t *testing.T) {
	w := New(0, 1*time.Second)
	w.Feed("anchor", ts(0), true)
	w.Feed("within", ts(1), false)
	out := w.Feed("outside", ts(5), false)
	for _, l := range out {
		if l == "outside" {
			t.Error("line outside after-window should not be emitted")
		}
	}
}

func TestReset_ClearsState(t *testing.T) {
	w := New(2*time.Second, 2*time.Second)
	w.Feed("a", ts(0), true)
	w.Reset()
	out := w.Feed("b", ts(1), false)
	if len(out) != 0 {
		t.Errorf("expected empty after reset, got %v", out)
	}
}

func TestFeed_MultipleAnchorsExtendWindow(t *testing.T) {
	w := New(0, 3*time.Second)
	w.Feed("a1", ts(0), true)
	w.Feed("mid", ts(2), false)
	w.Feed("a2", ts(4), true) // second anchor resets emitUntil
	out := w.Feed("late", ts(6), false)
	var found bool
	for _, l := range out {
		if l == "late" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'late' to be emitted after second anchor extended window")
	}
}
