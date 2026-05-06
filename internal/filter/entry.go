package filter

import "time"

// Entry represents a single parsed log line with structured fields.
type Entry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Raw       string
	Fields    map[string]string
}

// NewEntry creates an Entry with the given fields.
func NewEntry(ts time.Time, level, message, raw string) *Entry {
	return &Entry{
		Timestamp: ts,
		Level:     level,
		Message:   message,
		Raw:       raw,
		Fields:    make(map[string]string),
	}
}

// AddField adds a key-value pair to the entry's structured fields.
func (e *Entry) AddField(key, value string) {
	e.Fields[key] = value
}
