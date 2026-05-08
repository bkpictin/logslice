package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_FreshFile(t *testing.T) {
	cp, err := New(tmpPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := cp.Get("/var/log/app.log")
	if ok {
		t.Fatal("expected no state for unknown file")
	}
}

func TestSet_And_Get(t *testing.T) {
	cp, _ := New(tmpPath(t))
	cp.Set("/var/log/app.log", 1024, 42)
	s, ok := cp.Get("/var/log/app.log")
	if !ok {
		t.Fatal("expected state to be present")
	}
	if s.Offset != 1024 {
		t.Errorf("offset: got %d, want 1024", s.Offset)
	}
	if s.Inode != 42 {
		t.Errorf("inode: got %d, want 42", s.Inode)
	}
}

func TestFlush_And_Reload(t *testing.T) {
	p := tmpPath(t)
	cp, _ := New(p)
	cp.Set("/var/log/svc.log", 512, 7)
	if err := cp.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	cp2, err := New(p)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	s, ok := cp2.Get("/var/log/svc.log")
	if !ok {
		t.Fatal("expected reloaded state")
	}
	if s.Offset != 512 {
		t.Errorf("offset after reload: got %d, want 512", s.Offset)
	}
}

func TestReset_RemovesEntry(t *testing.T) {
	p := tmpPath(t)
	cp, _ := New(p)
	cp.Set("/var/log/app.log", 100, 1)
	cp.Reset("/var/log/app.log")
	_, ok := cp.Get("/var/log/app.log")
	if ok {
		t.Fatal("expected entry to be removed after Reset")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	if err := os.WriteFile(p, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for corrupt checkpoint file")
	}
}

func TestFlush_WritesAtomically(t *testing.T) {
	p := tmpPath(t)
	cp, _ := New(p)
	cp.Set("/a", 10, 0)
	if err := cp.Flush(); err != nil {
		t.Fatal(err)
	}
	// tmp file should not linger
	if _, err := os.Stat(p + ".tmp"); !os.IsNotExist(err) {
		t.Error("expected .tmp file to be cleaned up after atomic rename")
	}
}
