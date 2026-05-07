// Command logslice is a fast log file slicer and filter tool.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/pipeline"
)

func main() {
	var (
		start   = flag.String("start", "", "start time filter (RFC3339, e.g. 2024-01-01T00:00:00Z)")
		end     = flag.String("end", "", "end time filter (RFC3339)")
		level   = flag.String("level", "", "minimum log level (DEBUG, INFO, WARN, ERROR)")
		pattern = flag.String("pattern", "", "regex pattern to match against log lines")
		format  = flag.String("format", "text", "output format: text, json, csv")
		maxLines = flag.Int("max", 0, "maximum number of matching lines (0 = unlimited)")
		stats   = flag.Bool("stats", false, "print summary statistics to stderr after processing")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: logslice [options] [file ...]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf no files are given, reads from stdin.\n")
	}
	flag.Parse()

	cfg := pipeline.Config{
		Files:    flag.Args(),
		Format:   *format,
		Level:    *level,
		Pattern:  *pattern,
		MaxLines: *maxLines,
		Stats:    *stats,
	}

	if *start != "" {
		t, err := time.Parse(time.RFC3339, *start)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logslice: invalid --start: %v\n", err)
			os.Exit(1)
		}
		cfg.Start = t
	}
	if *end != "" {
		t, err := time.Parse(time.RFC3339, *end)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logslice: invalid --end: %v\n", err)
			os.Exit(1)
		}
		cfg.End = t
	}

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg pipeline.Config) error {
	if len(cfg.Files) == 0 {
		return pipeline.RunFromReader(os.Stdin, os.Stdout, cfg)
	}
	p, err := pipeline.New(cfg, os.Stdout)
	if err != nil {
		return err
	}
	return p.Run()
}
