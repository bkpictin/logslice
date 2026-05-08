package context_test

import (
	"context"
	"testing"
	"time"

	lsctx "github.com/yourorg/logslice/internal/context"
)

func TestWithStartTime_RoundTrip(t *testing.T) {
	now := time.Now()
	ctx := lsctx.WithStartTime(context.Background(), now)

	got, ok := lsctx.StartTime(ctx)
	if !ok {
		t.Fatal("expected start time to be present")
	}
	if !got.Equal(now) {
		t.Fatalf("expected %v, got %v", now, got)
	}
}

func TestStartTime_Missing(t *testing.T) {
	_, ok := lsctx.StartTime(context.Background())
	if ok {
		t.Fatal("expected no start time on empty context")
	}
}

func TestWithRequestID_RoundTrip(t *testing.T) {
	ctx := lsctx.WithRequestID(context.Background(), "run-99")

	got, ok := lsctx.RequestID(ctx)
	if !ok {
		t.Fatal("expected request ID to be present")
	}
	if got != "run-99" {
		t.Fatalf("expected run-99, got %q", got)
	}
}

func TestRequestID_Missing(t *testing.T) {
	_, ok := lsctx.RequestID(context.Background())
	if ok {
		t.Fatal("expected no request ID on empty context")
	}
}

func TestRequestID_EmptyString(t *testing.T) {
	ctx := lsctx.WithRequestID(context.Background(), "")
	_, ok := lsctx.RequestID(ctx)
	if ok {
		t.Fatal("empty string should be treated as absent")
	}
}

func TestElapsed_WithStartTime(t *testing.T) {
	ctx := lsctx.WithStartTime(context.Background(), time.Now().Add(-50*time.Millisecond))

	elapsed := lsctx.Elapsed(ctx)
	if elapsed < 50*time.Millisecond {
		t.Fatalf("expected elapsed >= 50ms, got %v", elapsed)
	}
}

func TestElapsed_NoStartTime(t *testing.T) {
	elapsed := lsctx.Elapsed(context.Background())
	if elapsed != 0 {
		t.Fatalf("expected 0 elapsed without start time, got %v", elapsed)
	}
}

func TestContexts_Independent(t *testing.T) {
	base := context.Background()
	ctx1 := lsctx.WithRequestID(base, "a")
	ctx2 := lsctx.WithRequestID(base, "b")

	id1, _ := lsctx.RequestID(ctx1)
	id2, _ := lsctx.RequestID(ctx2)

	if id1 == id2 {
		t.Fatal("contexts should be independent")
	}
}
