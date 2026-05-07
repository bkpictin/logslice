package pipeline

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/reader"
	"github.com/user/logslice/internal/stats"
)

// RunFromReader executes a pipeline reading log lines from r instead of files.
// This is the code path used when logslice reads from stdin.
func RunFromReader(ctx context.Context, cfg Config, r io.Reader, out io.Writer) (*stats.Stats, error) {
	f, err := filter.New(cfg.Level, cfg.Pattern, cfg.Since, cfg.Until)
	if err != nil {
		return nil, fmt.Errorf("pipeline: invalid filter: %w", err)
	}
	fmt, err := output.New(cfg.format(), out)
	if err != nil {
		return nil, fmt.Errorf("pipeline: invalid formatter: %w", err)
	}

	lr := reader.NewFromReader(r)
	st := stats.New()
	matched := 0

	for lr.Next() {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		line := lr.Line()
		st.AddLine(line)
		entry := filter.NewEntry(line)
		if !f.Match(entry) {
			continue
		}
		st.AddMatched(entry)
		if err := fmt.Write(entry); err != nil {
			return nil, fmt.Errorf("pipeline: write: %w", err)
		}
		matched++
		if cfg.MaxLines > 0 && matched >= cfg.MaxLines {
			break
		}
	}
	if err := lr.Err(); err != nil {
		return nil, fmt.Errorf("pipeline: read: %w", err)
	}
	st.Finish()
	_ = os.Stderr // suppress unused import if os is only used here
	return st, nil
}
