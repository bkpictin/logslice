package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
)

var testEntry = output.Entry{
	Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	Level:     "info",
	Message:   "application started",
}

func TestFormatter_Text(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	if err := f.Write(testEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %q", got)
	}
	if !strings.Contains(got, "application started") {
		t.Errorf("expected message in output, got: %q", got)
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)
	if err := f.Write(testEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, `"level":"info"`) {
		t.Errorf("expected level field in JSON, got: %q", got)
	}
	if !strings.Contains(got, `"message":"application started"`) {
		t.Errorf("expected message field in JSON, got: %q", got)
	}
}

func TestFormatter_CSV(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatCSV)
	if err := f.Write(testEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header + data), got %d", len(lines))
	}
	if lines[0] != "timestamp,level,message" {
		t.Errorf("unexpected CSV header: %q", lines[0])
	}
	if !strings.Contains(lines[1], "INFO") {
		t.Errorf("expected INFO in CSV row, got: %q", lines[1])
	}
}

func TestFormatter_CSV_HeaderOnce(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatCSV)
	for i := 0; i < 3; i++ {
		if err := f.Write(testEntry); err != nil {
			t.Fatalf("unexpected error on write %d: %v", i, err)
		}
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (1 header + 3 data), got %d", len(lines))
	}
}
