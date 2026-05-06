// Package output provides structured output formatting for log entries.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Format represents the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Entry represents a parsed log entry for output.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// Formatter writes log entries to an io.Writer in a given format.
type Formatter struct {
	format Format
	w      io.Writer
	first  bool
}

// New creates a new Formatter writing to w in the specified format.
func New(w io.Writer, format Format) *Formatter {
	return &Formatter{format: format, w: w, first: true}
}

// Write formats and writes a single log entry.
func (f *Formatter) Write(e Entry) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(e)
	case FormatCSV:
		return f.writeCSV(e)
	default:
		return f.writeText(e)
	}
}

func (f *Formatter) writeText(e Entry) error {
	_, err := fmt.Fprintf(f.w, "%s [%s] %s\n",
		e.Timestamp.Format(time.RFC3339), strings.ToUpper(e.Level), e.Message)
	return err
}

func (f *Formatter) writeJSON(e Entry) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f.w, "%s\n", data)
	return err
}

func (f *Formatter) writeCSV(e Entry) error {
	if f.first {
		_, err := fmt.Fprintln(f.w, "timestamp,level,message")
		if err != nil {
			return err
		}
		f.first = false
	}
	_, err := fmt.Fprintf(f.w, "%s,%s,%q\n",
		e.Timestamp.Format(time.RFC3339),
		strings.ToUpper(e.Level),
		e.Message)
	return err
}
