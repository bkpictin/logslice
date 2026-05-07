package pipeline_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/pipeline"
)

func TestRunTail_ReceivesLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	// Write initial content.
	_, _ = f.WriteString("first line\n")
	_ = f.Sync()

	var buf bytes.Buffer
	cfg := pipeline.Config{
		Files:  []string{f.Name()},
		Format: "text",
		Writer: &buf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Write a second line shortly after start.
	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = f.WriteString("second line\n")
		_ = f.Sync()
	}()

	err = pipeline.RunTail(ctx, cfg)
	if err != nil && err != context.DeadlineExceeded {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "second line") {
		t.Errorf("expected output to contain 'second line', got: %q", out)
	}
}

func TestRunTail_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	cfg := pipeline.Config{
		Files:  []string{"/nonexistent/path/file.log"},
		Format: "text",
		Writer: &buf,
	}

	err := pipeline.RunTail(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRunTail_ContextCancel(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-cancel-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString("line one\n")
	_ = f.Sync()

	var buf bytes.Buffer
	cfg := pipeline.Config{
		Files:  []string{f.Name()},
		Format: "text",
		Writer: &buf,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(40 * time.Millisecond)
		cancel()
	}()

	err = pipeline.RunTail(ctx, cfg)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
