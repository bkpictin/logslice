package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildStats(t *testing.T) *Stats {
	t.Helper()
	s := New()
	s.AddLine("error", true)
	s.AddLine("warn", true)
	s.AddLine("warn", false)
	s.AddLine("info", true)
	s.AddLine("", false)
	s.Finish()
	return s
}

func TestReport_ContainsExpectedFields(t *testing.T) {
	s := buildStats(t)
	var buf bytes.Buffer
	if err := s.Report(&buf); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Lines scanned:", "Lines matched:", "Duration:", "ERROR", "WARN", "INFO"} {
		if !strings.Contains(out, want) {
			t.Errorf("Report output missing %q\nGot:\n%s", want, out)
		}
	}
}

func TestReport_Counts(t *testing.T) {
	s := buildStats(t)
	var buf bytes.Buffer
	_ = s.Report(&buf)
	out := buf.String()
	if !strings.Contains(out, "5") {
		t.Errorf("expected total lines 5 in output, got:\n%s", out)
	}
}

func TestReportJSON_ValidOutput(t *testing.T) {
	s := buildStats(t)
	var buf bytes.Buffer
	if err := s.ReportJSON(&buf); err != nil {
		t.Fatalf("ReportJSON returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{`"total_lines":5`, `"matched_lines":3`, `"duration_ms":`} {
		if !strings.Contains(out, want) {
			t.Errorf("ReportJSON missing %q\nGot: %s", want, out)
		}
	}
}

func TestReportJSON_DurationPositive(t *testing.T) {
	s := New()
	time.Sleep(2 * time.Millisecond)
	s.Finish()
	var buf bytes.Buffer
	_ = s.ReportJSON(&buf)
	out := buf.String()
	if strings.Contains(out, `"duration_ms":0`) {
		t.Errorf("expected non-zero duration_ms, got: %s", out)
	}
}
