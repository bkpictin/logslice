package multiline_test

import (
	"testing"

	"github.com/logslice/logslice/internal/multiline"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := multiline.New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestFeed_SingleLineEntries(t *testing.T) {
	f, err := multiline.New(`^\d{4}-`)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// First entry starts — nothing flushed yet.
	_, ok := f.Feed("2024-01-01 INFO hello")
	if ok {
		t.Error("expected no flush on first entry")
	}

	// Second entry starts — first should be flushed.
	flushed, ok := f.Feed("2024-01-02 INFO world")
	if !ok {
		t.Error("expected flush when second entry begins")
	}
	if flushed != "2024-01-01 INFO hello" {
		t.Errorf("unexpected flushed value: %q", flushed)
	}
}

func TestFeed_ContinuationLinesAppended(t *testing.T) {
	f, err := multiline.New(`^ERROR`)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	f.Feed("ERROR something went wrong")
	f.Feed("\tat com.example.Foo.bar(Foo.java:42)")
	f.Feed("\tat com.example.Main.main(Main.java:10)")

	// Start a new entry to flush the previous.
	flushed, ok := f.Feed("ERROR another error")
	if !ok {
		t.Fatal("expected flush")
	}

	want := "ERROR something went wrong\n\tat com.example.Foo.bar(Foo.java:42)\n\tat com.example.Main.main(Main.java:10)"
	if flushed != want {
		t.Errorf("flushed = %q, want %q", flushed, want)
	}
}

func TestFlush_ReturnsBuffered(t *testing.T) {
	f, _ := multiline.New(`^INFO`)
	f.Feed("INFO line one")
	f.Feed("  continuation")

	got, ok := f.Flush()
	if !ok {
		t.Fatal("Flush should return true when data is buffered")
	}
	if got != "INFO line one\n  continuation" {
		t.Errorf("unexpected: %q", got)
	}

	// Second flush should be empty.
	_, ok = f.Flush()
	if ok {
		t.Error("second Flush should return false")
	}
}

func TestWithSeparator(t *testing.T) {
	f, _ := multiline.New(`^LOG`, multiline.WithSeparator(" | "))
	f.Feed("LOG start")
	f.Feed("cont1")
	f.Feed("cont2")

	got, _ := f.Flush()
	want := "LOG start | cont1 | cont2"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReset_DiscardsBuffer(t *testing.T) {
	f, _ := multiline.New(`^LOG`)
	f.Feed("LOG entry")
	f.Reset()
	_, ok := f.Flush()
	if ok {
		t.Error("expected Flush to return false after Reset")
	}
}
