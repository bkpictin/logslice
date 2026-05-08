// Package checkpoint tracks read progress through log files so that
// logslice can resume from where it left off across restarts.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// State holds the persisted position for a single file.
type State struct {
	File   string `json:"file"`
	Offset int64  `json:"offset"`
	Inode  uint64 `json:"inode,omitempty"`
}

// Checkpoint manages read-position persistence for one or more log files.
type Checkpoint struct {
	mu   sync.Mutex
	path string
	data map[string]State
}

// New loads an existing checkpoint file at path, or starts fresh if the file
// does not yet exist. Any other I/O error is returned.
func New(path string) (*Checkpoint, error) {
	c := &Checkpoint{
		path: path,
		data: make(map[string]State),
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return c, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(b, &c.data); err != nil {
		return nil, err
	}
	return c, nil
}

// Get returns the last saved State for file, and whether one was found.
func (c *Checkpoint) Get(file string) (State, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s, ok := c.data[file]
	return s, ok
}

// Set updates the in-memory position for file.
func (c *Checkpoint) Set(file string, offset int64, inode uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[file] = State{File: file, Offset: offset, Inode: inode}
}

// Flush writes the current state to disk atomically.
func (c *Checkpoint) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	b, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := c.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, c.path)
}

// Reset removes the saved position for file.
func (c *Checkpoint) Reset(file string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, file)
}
