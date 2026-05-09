// Package offset tracks byte offsets within log files, enabling
// resumable reads by persisting the last successfully processed
// position. It is safe for concurrent use.
package offset

import (
	"encoding/json"
	"os"
	"sync"
)

// Tracker records and persists byte offsets for one or more named files.
type Tracker struct {
	mu      sync.Mutex
	path    string
	offsets map[string]int64
}

// New creates a Tracker backed by the given state file.
// If the file exists its contents are loaded; a missing file is not an error.
func New(statePath string) (*Tracker, error) {
	t := &Tracker{
		path:    statePath,
		offsets: make(map[string]int64),
	}
	if err := t.load(); err != nil {
		return nil, err
	}
	return t, nil
}

// Get returns the last saved offset for the named file, or 0 if unknown.
func (t *Tracker) Get(name string) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.offsets[name]
}

// Set updates the in-memory offset for the named file.
func (t *Tracker) Set(name string, off int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if off < 0 {
		off = 0
	}
	t.offsets[name] = off
}

// Flush writes all current offsets to the backing state file.
func (t *Tracker) Flush() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.save()
}

// Reset removes the offset entry for the named file and flushes.
func (t *Tracker) Reset(name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.offsets, name)
	return t.save()
}

func (t *Tracker) load() error {
	data, err := os.ReadFile(t.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &t.offsets)
}

func (t *Tracker) save() error {
	data, err := json.Marshal(t.offsets)
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, data, 0o644)
}
