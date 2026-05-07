// Package tail provides a file-tailing utility for logslice.
//
// It watches a log file for new content appended after the process starts,
// emitting each new line on a channel. Polling is used for broad OS
// compatibility without requiring inotify or kqueue bindings.
//
// Basic usage:
//
//	t := tail.New("/var/log/app.log", 200*time.Millisecond)
//	go t.Run(ctx)
//	for line := range t.Lines() {
//		fmt.Println(line)
//	}
package tail
