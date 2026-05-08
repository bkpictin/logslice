package fieldextract_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldextract"
)

func TestExtract_EmptyLine(t *testing.T) {
	e := fieldextract.New()
	f := e.Extract("")
	if len(f) != 0 {
		t.Fatalf("expected empty fields, got %v", f)
	}
}

func TestExtract_JSON(t *testing.T) {
	e := fieldextract.New()
	f := e.Extract(`{"level":"info","msg":"hello","count":3}`)
	if f["level"] != "info" {
		t.Errorf("level: got %q, want %q", f["level"], "info")
	}
	if f["msg"] != "hello" {
		t.Errorf("msg: got %q, want %q", f["msg"], "hello")
	}
	if f["count"] != "3" {
		t.Errorf("count: got %q, want %q", f["count"], "3")
	}
}

func TestExtract_JSONWithLeadingTimestamp(t *testing.T) {
	e := fieldextract.New()
	line := `2024-01-02T15:04:05Z {"level":"warn","msg":"disk full"}`
	f := e.Extract(line)
	if f["level"] != "warn" {
		t.Errorf("level: got %q, want %q", f["level"], "warn")
	}
}

func TestExtract_KV(t *testing.T) {
	e := fieldextract.New()
	f := e.Extract(`ts=2024-01-02 level=error msg="disk full" retries=3`)
	if f["level"] != "error" {
		t.Errorf("level: got %q, want %q", f["level"], "error")
	}
	if f["msg"] != "disk full" {
		t.Errorf("msg: got %q, want %q", f["msg"], "disk full")
	}
	if f["retries"] != "3" {
		t.Errorf("retries: got %q, want %q", f["retries"], "3")
	}
}

func TestExtract_KVOnly(t *testing.T) {
	e := fieldextract.New(fieldextract.WithKV())
	// Disable JSON — a JSON line should fall through to KV.
	f := e.Extract(`key=val`)
	if f["key"] != "val" {
		t.Errorf("key: got %q, want %q", f["key"], "val")
	}
}

func TestExtract_JSONDisabled(t *testing.T) {
	e := fieldextract.New(fieldextract.WithKV())
	// Without JSON enabled, a pure JSON line yields no fields (no k=v pairs).
	f := e.Extract(`{"level":"info"}`)
	if _, ok := f["level"]; ok {
		t.Errorf("expected no level field when JSON disabled, got %q", f["level"])
	}
}

func TestExtract_NoFields(t *testing.T) {
	e := fieldextract.New()
	f := e.Extract("plain log line with no structure")
	if len(f) != 0 {
		t.Errorf("expected empty fields for unstructured line, got %v", f)
	}
}

func TestExtract_JSONPriorityOverKV(t *testing.T) {
	e := fieldextract.New()
	// Line contains both JSON and would match KV if JSON fails.
	line := `{"level":"debug","x":1}`
	f := e.Extract(line)
	if f["level"] != "debug" {
		t.Errorf("level: got %q, want %q", f["level"], "debug")
	}
}
