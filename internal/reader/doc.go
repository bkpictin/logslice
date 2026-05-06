// Package reader provides line-by-line reading of log files.
//
// It supports reading from regular files, stdin, and any io.Reader.
// Lines are emitted over a channel so that downstream processing
// (parsing, filtering, formatting) can be pipelined efficiently.
//
// Basic usage:
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
package reader
