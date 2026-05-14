package severity_test

import (
	"testing"

	"github.com/logslice/logslice/internal/levelmap"
	"github.com/logslice/logslice/internal/severity"
)

func TestNew_Disabled_WhenUnknown(t *testing.T) {
	f := severity.New(levelmap.Unknown)
	if f.Enabled() {
		t.Fatal("expected filter to be disabled for Unknown level")
	}
}

func TestNew_Enabled_WhenLevelSet(t *testing.T) {
	f := severity.New(levelmap.Warn)
	if !f.Enabled() {
		t.Fatal("expected filter to be enabled")
	}
}

func TestAllow_PassesWhenDisabled(t *testing.T) {
	f := severity.New(levelmap.Unknown)
	for _, lvl := range []string{"", "debug", "info", "warn", "error"} {
		if !f.Allow(lvl) {
			t.Fatalf("disabled filter should pass level %q", lvl)
		}
	}
}

func TestAllow_DropsBelow(t *testing.T) {
	f := severity.New(levelmap.Warn)
	if f.Allow("debug") {
		t.Error("debug should be dropped when min=warn")
	}
	if f.Allow("info") {
		t.Error("info should be dropped when min=warn")
	}
	if f.Dropped() != 2 {
		t.Errorf("expected 2 dropped, got %d", f.Dropped())
	}
}

func TestAllow_PassesAtOrAbove(t *testing.T) {
	f := severity.New(levelmap.Warn)
	for _, lvl := range []string{"warn", "error", "fatal"} {
		if !f.Allow(lvl) {
			t.Errorf("level %q should pass when min=warn", lvl)
		}
	}
	if f.Dropped() != 0 {
		t.Errorf("expected 0 dropped, got %d", f.Dropped())
	}
}

func TestReset_ClearsDropped(t *testing.T) {
	f := severity.New(levelmap.Error)
	f.Allow("debug")
	f.Allow("info")
	if f.Dropped() != 2 {
		t.Fatalf("expected 2 before reset, got %d", f.Dropped())
	}
	f.Reset()
	if f.Dropped() != 0 {
		t.Errorf("expected 0 after reset, got %d", f.Dropped())
	}
}

func TestAllow_EmptyLevel_DroppedWhenMinSet(t *testing.T) {
	f := severity.New(levelmap.Info)
	// empty string parses as Unknown which does not satisfy AtLeast(Info)
	if f.Allow("") {
		t.Error("empty level should be dropped when min=info")
	}
}
