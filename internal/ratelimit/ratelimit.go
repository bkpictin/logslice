// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how many log lines are emitted per second during pipeline
// processing. When the limit is zero or negative the limiter is a no-op.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which lines are passed through.
type Limiter struct {
	ticker  *time.Ticker
	disabled bool
}

// New creates a Limiter that allows at most linesPerSecond lines through per
// second. If linesPerSecond is <= 0 the limiter is disabled (no-op).
func New(linesPerSecond int) *Limiter {
	if linesPerSecond <= 0 {
		return &Limiter{disabled: true}
	}
	interval := time.Second / time.Duration(linesPerSecond)
	return &Limiter{
		ticker: time.NewTicker(interval),
	}
}

// Wait blocks until the next token is available or the context is cancelled.
// It returns ctx.Err() if the context is done before a token arrives, and nil
// when a token has been consumed.
func (l *Limiter) Wait(ctx context.Context) error {
	if l.disabled {
		return ctx.Err()
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.ticker.C:
		return nil
	}
}

// Stop releases resources held by the limiter. Safe to call on a disabled
// limiter.
func (l *Limiter) Stop() {
	if !l.disabled && l.ticker != nil {
		l.ticker.Stop()
	}
}

// Disabled reports whether the limiter is a no-op.
func (l *Limiter) Disabled() bool {
	return l.disabled
}
