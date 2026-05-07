package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/pipeline"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestRun_TextOutput(t *testing.T) {
	path := writeTempLog(t, []string{
		"2024-01-01T10:00:00Z INFO  hello world",
		"2024-01-01T10:01:00Z ERROR something failed",
	})
	cfg := pipeline.Config{
		Files:  []string{path},
		Format: "text",
	}
	if err := run(cfg); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestRun_MissingFile(t *testing.T) {
	cfg := pipeline.Config{
		Files:  []string{filepath.Join(t.TempDir(), "no-such-file.log")},
		Format: "text",
	}
	if err := run(cfg); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRun_StdinFallback(t *testing.T) {
	input := "2024-01-01T10:00:00Z INFO  stdin line\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = oldStdin })

	w.WriteString(input)
	w.Close()

	cfg := pipeline.Config{
		Files:  []string{},
		Format: "text",
	}
	if err := run(cfg); err != nil {
		t.Fatalf("run from stdin: %v", err)
	}
}

func TestRun_PatternFilter(t *testing.T) {
	path := writeTempLog(t, []string{
		"2024-01-01T10:00:00Z INFO  keep this line",
		"2024-01-01T10:01:00Z INFO  drop this one",
	})

	// Redirect stdout to capture output.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = old })

	cfg := pipeline.Config{
		Files:   []string{path},
		Format:  "text",
		Pattern: "keep",
	}
	if err := run(cfg); err != nil {
		t.Fatalf("run: %v", err)
	}
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "keep") {
		t.Errorf("expected 'keep' in output, got: %q", out)
	}
	if strings.Contains(out, "drop") {
		t.Errorf("unexpected 'drop' in output: %q", out)
	}
}
