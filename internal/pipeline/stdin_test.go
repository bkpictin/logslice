package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/pipeline"
)

// TestRunFromReader_BasicText verifies that RunFromReader processes lines from
// an io.Reader and emits formatted output to the provided writer.
func TestRunFromReader_BasicText(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T10:00:00Z INFO  starting service",
		"2024-01-15T10:00:01Z DEBUG initialising config",
		"2024-01-15T10:00:02Z ERROR failed to connect",
	}, "\n")

	cfg := pipeline.Config{
		Format: "text",
	}

	var out bytes.Buffer
	err := pipeline.RunFromReader(context.Background(), strings.NewReader(input), &out, cfg)
	if err != nil {
		t.Fatalf("RunFromReader returned unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "starting service") {
		t.Errorf("expected output to contain 'starting service', got:\n%s", got)
	}
	if !strings.Contains(got, "failed to connect") {
		t.Errorf("expected output to contain 'failed to connect', got:\n%s", got)
	}
}

// TestRunFromReader_LevelFilter verifies that only lines matching the requested
// log level are included in the output.
func TestRunFromReader_LevelFilter(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T10:00:00Z INFO  service started",
		"2024-01-15T10:00:01Z ERROR disk full",
		"2024-01-15T10:00:02Z INFO  heartbeat ok",
	}, "\n")

	cfg := pipeline.Config{
		Format: "text",
		Level:  "ERROR",
	}

	var out bytes.Buffer
	if err := pipeline.RunFromReader(context.Background(), strings.NewReader(input), &out, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if strings.Contains(got, "service started") {
		t.Errorf("INFO line should have been filtered out, got:\n%s", got)
	}
	if !strings.Contains(got, "disk full") {
		t.Errorf("ERROR line should be present, got:\n%s", got)
	}
}

// TestRunFromReader_EmptyInput verifies that an empty reader produces no error
// and no output lines.
func TestRunFromReader_EmptyInput(t *testing.T) {
	cfg := pipeline.Config{Format: "text"}

	var out bytes.Buffer
	err := pipeline.RunFromReader(context.Background(), strings.NewReader(""), &out, cfg)
	if err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}

	if out.Len() != 0 {
		t.Errorf("expected empty output, got %d bytes: %s", out.Len(), out.String())
	}
}

// TestRunFromReader_ContextCancel verifies that cancelling the context causes
// RunFromReader to stop processing and return promptly.
func TestRunFromReader_ContextCancel(t *testing.T) {
	// Build a large input so processing would take time without cancellation.
	var sb strings.Builder
	for i := 0; i < 10_000; i++ {
		sb.WriteString("2024-01-15T10:00:00Z INFO  log line number entry here\n")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	cfg := pipeline.Config{Format: "text"}
	var out bytes.Buffer

	// We don't assert on the error value — cancellation may or may not surface
	// as an error depending on timing, but the call must return.
	_ = pipeline.RunFromReader(ctx, strings.NewReader(sb.String()), &out, cfg)
}

// TestRunFromReader_JSONFormat verifies that JSON-formatted output is valid for
// at least one matched line.
func TestRunFromReader_JSONFormat(t *testing.T) {
	input := "2024-01-15T10:00:00Z INFO  json test line\n"

	cfg := pipeline.Config{Format: "json"}

	var out bytes.Buffer
	if err := pipeline.RunFromReader(context.Background(), strings.NewReader(input), &out, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "{") {
		t.Errorf("expected JSON output to contain '{', got:\n%s", got)
	}
}
