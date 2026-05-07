package tail

import (
	"context"
	"io"
	"os"
	"time"
)

// Tailer watches a file and emits new lines as they are appended.
type Tailer struct {
	path     string
	pollInterval time.Duration
	lines    chan string
}

// New creates a Tailer for the given file path.
// pollInterval controls how often the file is checked for new content.
func New(path string, pollInterval time.Duration) *Tailer {
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &Tailer{
		path:         path,
		pollInterval: pollInterval,
		lines:        make(chan string, 64),
	}
}

// Lines returns the channel on which new lines are delivered.
func (t *Tailer) Lines() <-chan string {
	return t.lines
}

// Run starts tailing the file, sending new lines to Lines().
// It blocks until ctx is cancelled or a read error occurs.
func (t *Tailer) Run(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only emit new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	buf := make([]byte, 0, 4096)
	ticker := time.NewTicker(t.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			chunk := make([]byte, 4096)
			n, readErr := f.Read(chunk)
			if n > 0 {
				buf = append(buf, chunk[:n]...)
				buf = t.flushLines(buf)
			}
			if readErr != nil && readErr != io.EOF {
				return readErr
			}
		}
	}
}

// flushLines extracts complete lines from buf, sends them, and returns remaining bytes.
func (t *Tailer) flushLines(buf []byte) []byte {
	start := 0
	for i, b := range buf {
		if b == '\n' {
			line := string(buf[start:i])
			if line != "" {
				t.lines <- line
			}
			start = i + 1
		}
	}
	return buf[start:]
}
