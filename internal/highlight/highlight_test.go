package highlight_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestLevelColor_KnownLevels(t *testing.T) {
	cases := []struct {
		level string
		want string
	}{
		{"ERROR", highlight.Red},
		{"FATAL", highlight.Red},
		{"WARN", highlight.Yellow},
		{"WARNING", highlight.Yellow},
		{"INFO", highlight.Green},
		{"DEBUG", highlight.Cyan},
		{"TRACE", highlight.Cyan},
		{"unknown", ""},
	}
	for _, tc := range cases {
		t.Run(tc.level, func(t *testing.T) {
			if got := highlight.LevelColor(tc.level); got != tc.want {
				t.Errorf("LevelColor(%q) = %q, want %q", tc.level, got, tc.want)
			}
		})
	}
}

func TestHighlighter_Disabled(t *testing.T) {
	h := highlight.New(false, nil)
	line := "2024-01-01 ERROR something bad"
	if got := h.Line(line, "ERROR"); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestHighlighter_Line_ColorWraps(t *testing.T) {
	h := highlight.New(true, nil)
	result := h.Line("some error", "ERROR")
	if !strings.HasPrefix(result, highlight.Red) {
		t.Error("expected line to start with Red escape")
	}
	if !strings.HasSuffix(result, highlight.Reset) {
		t.Error("expected line to end with Reset escape")
	}
}

func TestHighlighter_Line_PatternBolded(t *testing.T) {
	pat := regexp.MustCompile(`error`)
	h := highlight.New(true, pat)
	result := h.Line("an error occurred", "INFO")
	if !strings.Contains(result, highlight.Bold+"error"+highlight.Reset) {
		t.Errorf("expected bolded match in %q", result)
	}
}

func TestHighlighter_Word_Disabled(t *testing.T) {
	h := highlight.New(false, nil)
	if got := h.Word("ERROR", "ERROR"); got != "ERROR" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestHighlighter_Word_Colored(t *testing.T) {
	h := highlight.New(true, nil)
	result := h.Word("WARN", "WARN")
	if !strings.Contains(result, highlight.Yellow) {
		t.Errorf("expected yellow in %q", result)
	}
}
