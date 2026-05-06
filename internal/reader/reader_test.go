package reader_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func TestNewFromReader_EmitsLines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	lr := reader.NewFromReader(strings.NewReader(input), "test")
	defer lr.Close()

	var lines []reader.RawLine
	for l := range lr.Lines() {
		lines = append(lines, l)
	}

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Text != "line one" {
		t.Errorf("expected 'line one', got %q", lines[0].Text)
	}
	if lines[2].LineNum != 3 {
		t.Errorf("expected LineNum 3, got %d", lines[2].LineNum)
	}
}

func TestNewFromReader_EmptyInput(t *testing.T) {
	lr := reader.NewFromReader(strings.NewReader(""), "empty")
	defer lr.Close()

	var count int
	for range lr.Lines() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 lines for empty input, got %d", count)
	}
}

func TestNewFromReader_LineNumbers(t *testing.T) {
	input := "a\nb\nc\nd\n"
	lr := reader.NewFromReader(strings.NewReader(input), "nums")
	defer lr.Close()

	expected := 1
	for l := range lr.Lines() {
		if l.LineNum != expected {
			t.Errorf("line %d: expected LineNum %d, got %d", expected, expected, l.LineNum)
		}
		expected++
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := reader.New("/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestNewFromReader_ReadAtSet(t *testing.T) {
	lr := reader.NewFromReader(strings.NewReader("hello\n"), "ts")
	defer lr.Close()

	for l := range lr.Lines() {
		if l.ReadAt.IsZero() {
			t.Error("expected ReadAt to be set, got zero time")
		}
	}
}
