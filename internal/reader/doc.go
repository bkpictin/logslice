// Package reader provides line-by-line reading of log files.
//
// It supports reading from regular files, stdin, and any io.Reader.
// Lines are emitted over a channel so that downstream processing
// (parsing, filtering, formatting) can be pipelined efficiently.
//
// # Sources
//
// The reader can consume log lines from:
//   - Regular files on disk (via a file path)
//   - Standard input (os.Stdin)
//   - Any arbitrary [io.Reader] implementation
//
// # Basic usage
//
//	lr, err := reader.New("/var/log/app.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer lr.Close()
//
//	for line := range lr.Lines() {
//		fmt.Println(line.LineNum, line.Text)
//	}
//
// # Error handling
//
// If an error occurs during reading, the channel is closed and the error
// can be retrieved via [Reader.Err] after the channel has been drained.
package reader
