package offset_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/offset"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "offsets.json")
}

func TestNew_FreshFile(t *testing.T) {
	tr, err := offset.New(tmpPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tr.Get("app.log"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSet_And_Get(t *testing.T) {
	tr, _ := offset.New(tmpPath(t))
	tr.Set("app.log", 1024)
	if got := tr.Get("app.log"); got != 1024 {
		t.Fatalf("expected 1024, got %d", got)
	}
}

func TestSet_NegativeClampedToZero(t *testing.T) {
	tr, _ := offset.New(tmpPath(t))
	tr.Set("app.log", -99)
	if got := tr.Get("app.log"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestFlush_And_Reload(t *testing.T) {
	p := tmpPath(t)
	tr, _ := offset.New(p)
	tr.Set("a.log", 512)
	tr.Set("b.log", 2048)
	if err := tr.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	tr2, err := offset.New(p)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if got := tr2.Get("a.log"); got != 512 {
		t.Errorf("a.log: expected 512, got %d", got)
	}
	if got := tr2.Get("b.log"); got != 2048 {
		t.Errorf("b.log: expected 2048, got %d", got)
	}
}

func TestReset_RemovesEntry(t *testing.T) {
	p := tmpPath(t)
	tr, _ := offset.New(p)
	tr.Set("app.log", 999)
	_ = tr.Flush()

	if err := tr.Reset("app.log"); err != nil {
		t.Fatalf("reset error: %v", err)
	}
	if got := tr.Get("app.log"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}

	tr2, _ := offset.New(p)
	if got := tr2.Get("app.log"); got != 0 {
		t.Fatalf("expected 0 after reload, got %d", got)
	}
}

func TestNew_InvalidStateFile(t *testing.T) {
	p := tmpPath(t)
	if err := os.WriteFile(p, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := offset.New(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
