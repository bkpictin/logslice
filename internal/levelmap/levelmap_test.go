package levelmap_test

import (
	"testing"

	"github.com/logslice/logslice/internal/levelmap"
)

func TestParse_KnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  levelmap.Level
	}{
		{"trace", levelmap.Trace},
		{"TRC", levelmap.Trace},
		{"DEBUG", levelmap.Debug},
		{"dbg", levelmap.Debug},
		{"info", levelmap.Info},
		{"INF", levelmap.Info},
		{"warn", levelmap.Warn},
		{"WARNING", levelmap.Warn},
		{"WRN", levelmap.Warn},
		{"error", levelmap.Error},
		{"ERR", levelmap.Error},
		{"fatal", levelmap.Fatal},
		{"CRIT", levelmap.Fatal},
		{"critical", levelmap.Fatal},
		{"panic", levelmap.Fatal},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := levelmap.Parse(tc.input)
			if got != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParse_Unknown(t *testing.T) {
	for _, input := range []string{"", "verbose", "notice", "??"} {
		if got := levelmap.Parse(input); got != levelmap.Unknown {
			t.Errorf("Parse(%q) = %v, want Unknown", input, got)
		}
	}
}

func TestParse_Whitespace(t *testing.T) {
	if got := levelmap.Parse("  INFO  "); got != levelmap.Info {
		t.Errorf("expected Info, got %v", got)
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		level levelmap.Level
		want  string
	}{
		{levelmap.Unknown, "UNKNOWN"},
		{levelmap.Trace, "TRACE"},
		{levelmap.Debug, "DEBUG"},
		{levelmap.Info, "INFO"},
		{levelmap.Warn, "WARN"},
		{levelmap.Error, "ERROR"},
		{levelmap.Fatal, "FATAL"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestAtLeast(t *testing.T) {
	if !levelmap.AtLeast(levelmap.Error, levelmap.Warn) {
		t.Error("Error should be at least Warn")
	}
	if levelmap.AtLeast(levelmap.Debug, levelmap.Info) {
		t.Error("Debug should not be at least Info")
	}
	if !levelmap.AtLeast(levelmap.Info, levelmap.Info) {
		t.Error("Info should be at least Info")
	}
}
