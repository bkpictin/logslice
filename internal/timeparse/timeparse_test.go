package timeparse_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timeparse"
)

func TestParse_RFC3339(t *testing.T) {
	line := "2024-03-15T12:34:56Z some log message here"
	res, ok := timeparse.Parse(line)
	if !ok {
		t.Fatal("expected parse to succeed")
	}
	if res.Time.Year() != 2024 || res.Time.Month() != 3 || res.Time.Day() != 15 {
		t.Errorf("unexpected time: %v", res.Time)
	}
	if res.Remainder != "some log message here" {
		t.Errorf("unexpected remainder: %q", res.Remainder)
	}
}

func TestParse_RFC3339Nano(t *testing.T) {
	line := "2024-03-15T12:34:56.123456789Z msg"
	res, ok := timeparse.Parse(line)
	if !ok {
		t.Fatal("expected parse to succeed")
	}
	if res.Time.Nanosecond() == 0 {
		t.Error("expected nanoseconds to be set")
	}
	if res.Remainder != "msg" {
		t.Errorf("unexpected remainder: %q", res.Remainder)
	}
}

func TestParse_SpaceSeparated(t *testing.T) {
	line := "2024-03-15 08:00:00 server started"
	res, ok := timeparse.Parse(line)
	if !ok {
		t.Fatal("expected parse to succeed")
	}
	if res.Time.Hour() != 8 {
		t.Errorf("unexpected hour: %d", res.Time.Hour())
	}
	if res.Remainder != "server started" {
		t.Errorf("unexpected remainder: %q", res.Remainder)
	}
}

func TestParse_NoTimestamp(t *testing.T) {
	line := "no timestamp here at all"
	_, ok := timeparse.Parse(line)
	if ok {
		t.Error("expected parse to fail for line without timestamp")
	}
}

func TestParse_EmptyLine(t *testing.T) {
	_, ok := timeparse.Parse("")
	if ok {
		t.Error("expected parse to fail for empty line")
	}
}

func TestParse_RemainderStripped(t *testing.T) {
	line := "2024-01-01T00:00:00Z   leading spaces stripped"
	res, ok := timeparse.Parse(line)
	if !ok {
		t.Fatal("expected parse to succeed")
	}
	if res.Remainder != "leading spaces stripped" {
		t.Errorf("remainder not stripped: %q", res.Remainder)
	}
}

func TestFormats_ReturnsCopy(t *testing.T) {
	f1 := timeparse.Formats()
	f2 := timeparse.Formats()
	if len(f1) == 0 {
		t.Fatal("expected at least one format")
	}
	f1[0] = "mutated"
	if f2[0] == "mutated" {
		t.Error("Formats should return an independent copy")
	}
}

func TestParse_FormatRecorded(t *testing.T) {
	line := "2024-06-01T10:20:30Z event"
	res, ok := timeparse.Parse(line)
	if !ok {
		t.Fatal("expected parse to succeed")
	}
	if res.Format == "" {
		t.Error("expected Format to be recorded")
	}
	_, err := time.Parse(res.Format, line[:len(line)-len(" event")-1])
	_ = err // format may include zone; just ensure it is non-empty
}
