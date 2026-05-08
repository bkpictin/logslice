// Package rotate provides log file rotation detection for logslice.
// It watches a file path and signals when the file has been rotated
// (i.e. replaced or truncated), allowing readers to reopen it.
package rotate

import (
	"os"
	"time"
)

// Detector watches a file for rotation events such as truncation or
// replacement (inode change). It polls at a configurable interval.
type Detector struct {
	path     string
	interval time.Duration
	inode    uint64
	size     int64
}

// New creates a Detector for the given file path, polling at interval.
// It records the initial inode and size of the file.
func New(path string, interval time.Duration) (*Detector, error) {
	d := &Detector{path: path, interval: interval}
	if err := d.snapshot(); err != nil {
		return nil, err
	}
	return d, nil
}

// Rotated returns true if the file has been rotated since the last call
// to Reset (or construction). A rotation is detected when the inode
// changes or the file size is smaller than previously observed.
func (d *Detector) Rotated() (bool, error) {
	fi, err := os.Stat(d.path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	currentInode := inode(fi)
	currentSize := fi.Size()
	if currentInode != d.inode || currentSize < d.size {
		return true, nil
	}
	d.size = currentSize
	return false, nil
}

// Reset re-snapshots the current file state, treating it as the new
// baseline. Call this after reopening a rotated file.
func (d *Detector) Reset() error {
	return d.snapshot()
}

// Interval returns the polling interval configured for this detector.
func (d *Detector) Interval() time.Duration {
	return d.interval
}

func (d *Detector) snapshot() error {
	fi, err := os.Stat(d.path)
	if err != nil {
		return err
	}
	d.inode = inode(fi)
	d.size = fi.Size()
	return nil
}
