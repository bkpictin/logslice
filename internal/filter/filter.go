package filter

import (
	"regexp"
	"strings"
	"time"
)

// Options holds all filtering criteria for log lines.
type Options struct {
	StartTime  *time.Time
	EndTime    *time.Time
	Pattern    string
	Regexp     *regexp.Regexp
	Level      string
	CaseSensitive bool
}

// Filter applies filter options to a parsed log line and returns true if the line matches.
type Filter struct {
	opts Options
}

// New creates a new Filter with the given options.
func New(opts Options) (*Filter, error) {
	if opts.Pattern != "" && opts.Regexp == nil {
		flags := "(?i)"
		if opts.CaseSensitive {
			flags = ""
		}
		re, err := regexp.Compile(flags + opts.Pattern)
		if err != nil {
			return nil, err
		}
		opts.Regexp = re
	}
	return &Filter{opts: opts}, nil
}

// Match returns true if the given log entry matches all active filter criteria.
func (f *Filter) Match(entry *Entry) bool {
	if f.opts.StartTime != nil && entry.Timestamp.Before(*f.opts.StartTime) {
		return false
	}
	if f.opts.EndTime != nil && entry.Timestamp.After(*f.opts.EndTime) {
		return false
	}
	if f.opts.Level != "" {
		if !strings.EqualFold(entry.Level, f.opts.Level) {
			return false
		}
	}
	if f.opts.Regexp != nil {
		if !f.opts.Regexp.MatchString(entry.Message) {
			return false
		}
	}
	return true
}
