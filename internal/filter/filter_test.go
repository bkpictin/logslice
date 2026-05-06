package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func baseEntry() *filter.Entry {
	ts := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	return filter.NewEntry(ts, "ERROR", "connection refused on port 5432", "raw line")
}

func TestMatch_NoFilters(t *testing.T) {
	f, err := filter.New(filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(baseEntry()) {
		t.Error("expected entry to match with no filters")
	}
}

func TestMatch_TimeRange(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	f, _ := filter.New(filter.Options{StartTime: &start, EndTime: &end})
	if !f.Match(baseEntry()) {
		t.Error("expected entry within time range to match")
	}

	before := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	f2, _ := filter.New(filter.Options{StartTime: &before, EndTime: &before})
	if f2.Match(baseEntry()) {
		t.Error("expected entry outside time range to not match")
	}
}

func TestMatch_Level(t *testing.T) {
	f, _ := filter.New(filter.Options{Level: "error"})
	if !f.Match(baseEntry()) {
		t.Error("expected case-insensitive level match")
	}

	f2, _ := filter.New(filter.Options{Level: "INFO"})
	if f2.Match(baseEntry()) {
		t.Error("expected non-matching level to fail")
	}
}

func TestMatch_Pattern(t *testing.T) {
	f, err := filter.New(filter.Options{Pattern: "port \\d+"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(baseEntry()) {
		t.Error("expected pattern to match message")
	}
}

func TestMatch_InvalidPattern(t *testing.T) {
	_, err := filter.New(filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestMatch_CaseSensitivePattern(t *testing.T) {
	f, _ := filter.New(filter.Options{Pattern: "CONNECTION", CaseSensitive: true})
	if f.Match(baseEntry()) {
		t.Error("expected case-sensitive pattern to not match lowercase message")
	}
}
