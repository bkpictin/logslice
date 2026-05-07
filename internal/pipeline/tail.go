package pipeline

import (
	"context"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/stats"
	"github.com/user/logslice/internal/tail"
)

// RunTail follows a file for new log lines and processes them through the
// filter, formatter, and stats pipeline until the context is cancelled.
func RunTail(ctx context.Context, cfg Config) error {
	f, err := filter.New(cfg.StartTime, cfg.EndTime, cfg.Level, cfg.Pattern)
	if err != nil {
		return err
	}

	fmt, err := output.New(cfg.Format, cfg.Writer)
	if err != nil {
		return err
	}

	st := stats.New()

	t, err := tail.New(cfg.Files[0])
	if err != nil {
		return err
	}
	defer t.Close()

	lines := t.Lines()

	for {
		select {
		case <-ctx.Done():
			st.Finish()
			if cfg.ShowStats {
				st.Report(cfg.Writer)
			}
			return ctx.Err()

		case line, ok := <-lines:
			if !ok {
				st.Finish()
				if cfg.ShowStats {
					st.Report(cfg.Writer)
				}
				return nil
			}

			entry := filter.NewEntry(line.Text, line.Number)
			st.AddLine(entry)

			if !f.Match(entry) {
				continue
			}

			st.AddMatched(entry)

			if err := fmt.Write(entry); err != nil {
				return err
			}
		}
	}
}
