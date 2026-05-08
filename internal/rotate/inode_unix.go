//go:build !windows

package rotate

import (
	"os"
	"syscall"
)

// inode returns the inode number for the given FileInfo on Unix systems.
func inode(fi os.FileInfo) uint64 {
	if sys, ok := fi.Sys().(*syscall.Stat_t); ok {
		return sys.Ino
	}
	return 0
}
