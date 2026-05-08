package pipeline

import (
	"context"
	"fmt"
	"time"

	lsctx "github.com/yourorg/logslice/internal/context"
)

// enrichContext attaches pipeline-level metadata to the given context.
// It always sets a start time; if cfg provides a run ID it is set too.
func enrichContext(ctx context.Context, cfg *Config) context.Context {
	ctx = lsctx.WithStartTime(ctx, time.Now())
	if id := runID(cfg); id != "" {
		ctx = lsctx.WithRequestID(ctx, id)
	}
	return ctx
}

// runID derives an optional run identifier from the config.
// It uses the first input file name when available.
func runID(cfg *Config) string {
	if cfg == nil {
		return ""
	}
	if len(cfg.Files) > 0 && cfg.Files[0] != "" {
		return fmt.Sprintf("file:%s", cfg.Files[0])
	}
	return ""
}

// logContextInfo writes a short diagnostic line about the run context to
// stderr when the context carries a request ID. This is a no-op in normal
// production use but aids debugging.
func logContextInfo(ctx context.Context) {
	if id, ok := lsctx.RequestID(ctx); ok {
		_ = id // available for structured logging hooks
	}
}
