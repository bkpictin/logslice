package window

import "time"

// Window holds a sliding time window of log lines, emitting lines that fall
// within [start, end) relative to the first matched anchor line.
type Window struct {
	before   time.Duration
	after    time.Duration
	buf      []Entry
	anchor   *time.Time
	emitting bool
	emitUntil time.Time
}

// Entry is a timestamped log line held in the window buffer.
type Entry struct {
	Line string
	At   time.Time
}

// New creates a Window that captures `before` duration before an anchor match
// and `after` duration following it.
func New(before, after time.Duration) *Window {
	return &Window{
		before: before,
		after:  after,
		buf:    make([]Entry, 0, 64),
	}
}

// Feed adds a line with its timestamp. If match is true the line is treated as
// an anchor that opens an emission window. Feed returns all lines that should
// be emitted now.
func (w *Window) Feed(line string, at time.Time, match bool) []string {
	entry := Entry{Line: line, At: at}

	// Prune buffer entries that are too old to ever be emitted.
	if !w.emitting {
		cutoff := at.Add(-w.before)
		i := 0
		for i < len(w.buf) && w.buf[i].At.Before(cutoff) {
			i++
		}
		w.buf = w.buf[i:]
	}

	w.buf = append(w.buf, entry)

	if match {
		w.emitting = true
		w.emitUntil = at.Add(w.after)
	}

	if !w.emitting {
		return nil
	}

	var out []string
	remaining := w.buf[:0]
	for _, e := range w.buf {
		if !e.At.After(w.emitUntil) {
			out = append(out, e.Line)
		} else {
			remaining = append(remaining, e)
		}
	}
	w.buf = remaining

	if at.After(w.emitUntil) {
		w.emitting = false
	}
	return out
}

// Reset clears all buffered state.
func (w *Window) Reset() {
	w.buf = w.buf[:0]
	w.anchor = nil
	w.emitting = false
}
