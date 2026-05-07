package stats_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/stats"
)

func TestNew_InitialState(t *testing.T) {
	c := stats.New()
	if c.TotalLines != 0 || c.MatchedLines != 0 || c.SkippedLines != 0 {
		t.Errorf("expected zero counts, got total=%d matched=%d skipped=%d",
			c.TotalLines, c.MatchedLines, c.SkippedLines)
	}
	if c.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}
}

func TestAddLine_Matched(t *testing.T) {
	c := stats.New()
	c.AddLine(true, "ERROR", 42)
	c.AddLine(true, "WARN", 10)
	c.AddLine(false, "", 8)

	if c.TotalLines != 3 {
		t.Errorf("TotalLines: want 3, got %d", c.TotalLines)
	}
	if c.MatchedLines != 2 {
		t.Errorf("MatchedLines: want 2, got %d", c.MatchedLines)
	}
	if c.SkippedLines != 1 {
		t.Errorf("SkippedLines: want 1, got %d", c.SkippedLines)
	}
	if c.BytesRead != 60 {
		t.Errorf("BytesRead: want 60, got %d", c.BytesRead)
	}
}

func TestLevelCounts(t *testing.T) {
	c := stats.New()
	c.AddLine(true, "ERROR", 1)
	c.AddLine(true, "ERROR", 1)
	c.AddLine(true, "INFO", 1)
	c.AddLine(false, "DEBUG", 1) // skipped — should not appear

	lc := c.LevelCounts()
	if lc["ERROR"] != 2 {
		t.Errorf("ERROR count: want 2, got %d", lc["ERROR"])
	}
	if lc["INFO"] != 1 {
		t.Errorf("INFO count: want 1, got %d", lc["INFO"])
	}
	if _, ok := lc["DEBUG"]; ok {
		t.Error("DEBUG should not appear in level counts for skipped lines")
	}
}

func TestFinish_SetsEndTime(t *testing.T) {
	c := stats.New()
	time.Sleep(2 * time.Millisecond)
	c.Finish()

	if c.EndTime.IsZero() {
		t.Error("expected EndTime to be set after Finish")
	}
	if c.Duration() < time.Millisecond {
		t.Errorf("expected duration >= 1ms, got %s", c.Duration())
	}
}

func TestDuration_BeforeFinish(t *testing.T) {
	c := stats.New()
	time.Sleep(1 * time.Millisecond)
	d := c.Duration()
	if d <= 0 {
		t.Errorf("expected positive duration before Finish, got %s", d)
	}
}
