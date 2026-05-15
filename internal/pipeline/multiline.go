package pipeline

import (
	"context"

	"github.com/logslice/logslice/internal/multiline"
	"github.com/logslice/logslice/internal/reader"
)

// newMultilineCollector constructs a multiline.Collector from the pipeline
// Config. Returns nil when no multiline start pattern is configured.
func newMultilineCollector(cfg Config) (*multiline.Collector, error) {
	if cfg.MultilineStart == "" {
		return nil, nil
	}
	opts := []multiline.Option{}
	if cfg.MultilineSeparator != "" {
		opts = append(opts, multiline.WithSeparator(cfg.MultilineSeparator))
	}
	return multiline.New(cfg.MultilineStart, opts...)
}

// applyMultiline wraps the raw line channel produced by a reader with
// multiline folding. When collector is nil the original channel is returned
// unchanged so callers never need to branch.
func applyMultiline(
	ctx context.Context,
	lines <-chan reader.Line,
	collector *multiline.Collector,
) <-chan reader.Line {
	if collector == nil {
		return lines
	}

	out := make(chan reader.Line)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				// Flush any buffered partial entry before exiting.
				if entry, ok := collector.Flush(); ok {
					select {
					case out <- reader.Line{Text: entry}:
					case <-ctx.Done():
					}
				}
				return
			case line, ok := <-lines:
				if !ok {
					// Channel closed — flush remainder.
					if entry, fok := collector.Flush(); fok {
						select {
						case out <- reader.Line{Text: entry, Number: line.Number}:
						case <-ctx.Done():
						}
					}
					return
				}
				if entry, emitted := collector.Feed(line.Text); emitted {
					select {
					case out <- reader.Line{Text: entry, Number: line.Number}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return out
}
