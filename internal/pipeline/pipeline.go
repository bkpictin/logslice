package pipeline

import (
	"context"
	"fmt"
	"io"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/reader"
	"github.com/user/logslice/internal/stats"
)

// Pipeline orchestrates reading, filtering, formatting and stats collection.
type Pipeline struct {
	cfg    Config
	filter *filter.Filter
	fmt    *output.Formatter
	st     *stats.Stats
	out    io.Writer
}

// New constructs a Pipeline from cfg, writing output to out.
func New(cfg Config, out io.Writer) (*Pipeline, error) {
	f, err := filter.New(cfg.Level, cfg.Pattern, cfg.Since, cfg.Until)
	if err != nil {
		return nil, fmt.Errorf("pipeline: invalid filter: %w", err)
	}
	fmt, err := output.New(cfg.format(), out)
	if err != nil {
		return nil, fmt.Errorf("pipeline: invalid formatter: %w", err)
	}
	return &Pipeline{
		cfg:    cfg,
		filter: f,
		fmt:    fmt,
		st:     stats.New(),
		out:    out,
	}, nil
}

// Run executes the pipeline and returns collected statistics.
// It honours ctx cancellation between lines.
func (p *Pipeline) Run(ctx context.Context) (*stats.Stats, error) {
	var r *reader.Reader
	var err error

	if len(p.cfg.Files) == 0 {
		return nil, fmt.Errorf("pipeline: no input files specified")
	}
	if len(p.cfg.Files) == 1 {
		r, err = reader.New(p.cfg.Files[0])
	} else {
		r, err = reader.NewMulti(p.cfg.Files)
	}
	if err != nil {
		return nil, fmt.Errorf("pipeline: open input: %w", err)
	}
	defer r.Close()

	matched := 0
	for r.Next() {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		line := r.Line()
		p.st.AddLine(line)
		entry := filter.NewEntry(line)
		if !p.filter.Match(entry) {
			continue
		}
		p.st.AddMatched(entry)
		if err := p.fmt.Write(entry); err != nil {
			return nil, fmt.Errorf("pipeline: write: %w", err)
		}
		matched++
		if p.cfg.MaxLines > 0 && matched >= p.cfg.MaxLines {
			break
		}
	}
	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("pipeline: read: %w", err)
	}
	p.st.Finish()
	return p.st, nil
}
