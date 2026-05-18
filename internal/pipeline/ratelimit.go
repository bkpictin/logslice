package pipeline

import (
	"context"

	"github.com/logslice/logslice/internal/ratelimit"
)

// newRateLimiter constructs a ratelimit.Limiter from the pipeline config.
// Returns nil when LinesPerSecond is zero or negative (disabled).
func newRateLimiter(cfg Config) (*ratelimit.Limiter, error) {
	if cfg.LinesPerSecond <= 0 {
		return nil, nil
	}
	return ratelimit.New(cfg.LinesPerSecond), nil
}

// applyRateLimit blocks until the limiter allows the next line to proceed.
// If the limiter is nil (disabled) the line is returned immediately.
// A cancelled context causes the function to return an empty string and the
// context error so the caller can abort the pipeline.
func applyRateLimit(ctx context.Context, rl *ratelimit.Limiter, line string) (string, error) {
	if rl == nil {
		return line, nil
	}
	if err := rl.Wait(ctx); err != nil {
		return "", err
	}
	return line, nil
}
