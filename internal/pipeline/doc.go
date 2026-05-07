// Package pipeline wires together the reader, filter, stats, and output
// components into a single processing pipeline for logslice.
//
// A Pipeline reads log lines from one or more sources, applies filters,
// tracks statistics, and writes formatted output to a destination writer.
//
// Basic usage:
//
//	cfg := pipeline.Config{
//		Files:      []string{"app.log"},
//		Format:     "json",
//		Level:      "error",
//		Pattern:    "timeout",
//	}
//	p, err := pipeline.New(cfg, os.Stdout)
//	if err != nil {
//		log.Fatal(err)
//	}
//	stats, err := p.Run(context.Background())
package pipeline
