// Package truncate provides line-length limiting for log output.
//
// When processing high-volume or noisy log files it is useful to cap the
// length of individual lines so that wide terminal windows or downstream
// parsers are not overwhelmed by unexpectedly long entries.
//
// Basic usage:
//
//	tr := truncate.NewDefault(120)
//	fmt.Println(tr.Line(veryLongLogLine))
//
// Truncation is transparent when maxLen <= 0, making it safe to wire into
// the pipeline unconditionally and control behaviour via configuration.
package truncate
