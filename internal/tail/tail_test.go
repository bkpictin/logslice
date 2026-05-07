package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func writeLater(t *testing.T, f *os.File, lines []string, delay time.Duration) {
	t.Helper()
	time.Sleep(delay)
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tlr := tail.New(f.Name(), 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go writeLater(t, f, []string{"hello world", "second line"}, 100*time.Millisecond)

	go tlr.Run(ctx) //nolint:errcheck

	var got []string
	for len(got) < 2 {
		select {
		case line := <-tlr.Lines():
			got = append(got, line)
		case <-ctx.Done():
			t.Fatalf("timed out waiting for lines, got %v", got)
		}
	}

	if got[0] != "hello world" {
		t.Errorf("expected 'hello world', got %q", got[0])
	}
	if got[1] != "second line" {
		t.Errorf("expected 'second line', got %q", got[1])
	}
}

func TestTailer_MissingFile(t *testing.T) {
	tlr := tail.New("/nonexistent/path/file.log", 50*time.Millisecond)
	ctx := context.Background()
	err := tlr.Run(ctx)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestTailer_ContextCancel(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-cancel-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tlr := tail.New(f.Name(), 50*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- tlr.Run(ctx) }()

	time.Sleep(120 * time.Millisecond)
	cancel()

	select {
	case runErr := <-done:
		if runErr != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", runErr)
		}
	case <-time.After(2 * time.Second):
		t.Error("Run did not return after context cancel")
	}
}
