package stats

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// Report writes a human-readable summary of the collected statistics to w.
func (s *Stats) Report(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Lines scanned:\t%d\n", s.TotalLines())
	fmt.Fprintf(tw, "Lines matched:\t%d\n", s.MatchedLines())
	fmt.Fprintf(tw, "Duration:\t%s\n", s.Duration().Round(1000000))

	if len(s.LevelCounts()) > 0 {
		fmt.Fprintf(tw, "\nLevel breakdown:\n")

		levels := make([]string, 0, len(s.LevelCounts()))
		for lvl := range s.LevelCounts() {
			levels = append(levels, lvl)
		}
		sort.Strings(levels)

		for _, lvl := range levels {
			label := strings.ToUpper(lvl)
			if label == "" {
				label = "(unknown)"
			}
			fmt.Fprintf(tw, "  %s:\t%d\n", label, s.LevelCounts()[lvl])
		}
	}

	return tw.Flush()
}

// ReportJSON writes a JSON summary of statistics to w.
func (s *Stats) ReportJSON(w io.Writer) error {
	counts := s.LevelCounts()
	parts := make([]string, 0, len(counts))
	for lvl, n := range counts {
		key := lvl
		if key == "" {
			key = "unknown"
		}
		parts = append(parts, fmt.Sprintf(`%q:%d`, key, n))
	}
	sort.Strings(parts)

	levelJSON := "{}" 
	if len(parts) > 0 {
		levelJSON = "{" + strings.Join(parts, ",") + "}"
	}

	_, err := fmt.Fprintf(w,
		`{"total_lines":%d,"matched_lines":%d,"duration_ms":%d,"levels":%s}\n`,
		s.TotalLines(),
		s.MatchedLines(),
		s.Duration().Milliseconds(),
		levelJSON,
	)
	return err
}
