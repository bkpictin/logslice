// Package checkpoint provides durable read-position tracking for log files.
//
// A Checkpoint records the byte offset (and optional inode) reached for each
// file being processed. On the next run the pipeline can seek directly to the
// saved offset, skipping lines that were already handled.
//
// Usage:
//
//	cp, err := checkpoint.New("/var/run/logslice/state.json")
//	if err != nil { ... }
//
//	// On startup, resume from last position:
//	if s, ok := cp.Get(filename); ok {
//		file.Seek(s.Offset, io.SeekStart)
//	}
//
//	// After each batch, persist progress:
//	cp.Set(filename, currentOffset, currentInode)
//	cp.Flush()
package checkpoint
