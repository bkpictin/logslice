package tail_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

// TestTailer_HighVolume verifies the tailer handles many lines without dropping.
func TestTailer_HighVolume(t *testing.T) {
	const lineCount = 200

	f, err := os.CreateTemp(t.TempDir(), "tail-hv-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tlr := tail.New(f.Name(), 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go tlr.Run(ctx) //nolint:errcheck

	// Give the tailer a moment to seek to end before writing.
	time.Sleep(60 * time.Millisecond)

	go func() {
		for i := 0; i < lineCount; i++ {
			fmt.Fprintf(f, "line %d\n", i)
		}
	}()

	received := 0
	for received < lineCount {
		select {
		case <-tlr.Lines():
			received++
		case <-ctx.Done():
			t.Fatalf("timed out: only received %d/%d lines", received, lineCount)
		}
	}

	if received != lineCount {
		t.Errorf("expected %d lines, got %d", lineCount, received)
	}
}
