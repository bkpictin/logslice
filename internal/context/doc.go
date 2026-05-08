// Package context provides helpers for propagating logslice pipeline metadata
// through standard context.Context values.
//
// It stores and retrieves well-typed values such as pipeline start times and
// optional run identifiers without exposing raw context keys to callers.
//
// Usage:
//
//	ctx = lsctx.WithStartTime(ctx, time.Now())
//	ctx = lsctx.WithRequestID(ctx, "run-42")
//
//	elapsed := lsctx.Elapsed(ctx)
//	id, ok := lsctx.RequestID(ctx)
package context
