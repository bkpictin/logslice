package reader

import (
	"bufio"
	"io"
	"os"
	"time"
)

// LineReader reads log lines from a file or reader and emits RawLine values.
type LineReader struct {
	source  io.ReadCloser
	scanner *bufio.Scanner
}

// RawLine represents a single line read from a log source.
type RawLine struct {
	Text      string
	LineNum   int
	ReadAt    time.Time
	SourceName string
}

// New creates a LineReader from the given file path.
// Pass "-" to read from stdin.
func New(path string) (*LineReader, error) {
	var rc io.ReadCloser
	if path == "-" {
		rc = io.NopCloser(os.Stdin)
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		rc = f
	}
	return &LineReader{
		source:  rc,
		scanner: bufio.NewScanner(rc),
	}, nil
}

// NewFromReader creates a LineReader from an existing io.Reader.
func NewFromReader(r io.Reader, sourceName string) *LineReader {
	return &LineReader{
		source:  io.NopCloser(r),
		scanner: bufio.NewScanner(r),
	}
}

// Lines returns a channel that emits each RawLine in order.
// The channel is closed when EOF is reached or an error occurs.
func (lr *LineReader) Lines() <-chan RawLine {
	ch := make(chan RawLine, 64)
	go func() {
		defer close(ch)
		lineNum := 0
		for lr.scanner.Scan() {
			lineNum++
			ch <- RawLine{
				Text:    lr.scanner.Text(),
				LineNum: lineNum,
				ReadAt:  time.Now(),
			}
		}
	}()
	return ch
}

// Close releases the underlying resource.
func (lr *LineReader) Close() error {
	return lr.source.Close()
}
