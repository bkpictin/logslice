package pipeline_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/logslice/internal/pipeline"
)

func writeTempLog(t *testing.T, lines string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(lines); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestPipeline_RunText(t *testing.T) {
	path := writeTempLog(t, "line one\nline two\nline three\n")
	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{Files: []string{path}}, &buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	st, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if st.TotalLines() != 3 {
		t.Errorf("TotalLines = %d, want 3", st.TotalLines())
	}
	if !strings.Contains(buf.String(), "line one") {
		t.Errorf("output missing expected content: %q", buf.String())
	}
}

func TestPipeline_MaxLines(t *testing.T) {
	path := writeTempLog(t, "a\nb\nc\nd\ne\n")
	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{Files: []string{path}, MaxLines: 2}, &buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	st, err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if st.MatchedLines() != 2 {
		t.Errorf("MatchedLines = %d, want 2", st.MatchedLines())
	}
}

func TestPipeline_ContextCancel(t *testing.T) {
	path := writeTempLog(t, "x\ny\nz\n")
	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{Files: []string{path}}, &buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = p.Run(ctx)
	if err == nil {
		t.Error("expected error on cancelled context, got nil")
	}
}

func TestPipeline_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{
		Files: []string{filepath.Join(t.TempDir(), "no-such.log")},
	}, &buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, err = p.Run(context.Background())
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
