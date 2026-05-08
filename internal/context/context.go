package context

import (
	"context"
	"time"
)

// Key is an unexported type for context keys in this package.
type Key int

const (
	// startTimeKey stores the pipeline start time.
	startTimeKey Key = iota
	// requestIDKey stores an optional request/run identifier.
	requestIDKey
)

// WithStartTime returns a new context carrying the given start time.
func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, startTimeKey, t)
}

// StartTime retrieves the start time from ctx.
// If not set, it returns the zero time and false.
func StartTime(ctx context.Context) (time.Time, bool) {
	v, ok := ctx.Value(startTimeKey).(time.Time)
	return v, ok
}

// WithRequestID returns a new context carrying the given request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestID retrieves the request ID from ctx.
// If not set, it returns an empty string and false.
func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok && v != ""
}

// Elapsed returns the duration since the start time stored in ctx.
// If no start time is set, it returns 0.
func Elapsed(ctx context.Context) time.Duration {
	if t, ok := StartTime(ctx); ok {
		return time.Since(t)
	}
	return 0
}
