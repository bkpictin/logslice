package sample_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/sample"
)

func TestNew_RateOne_KeepsAll(t *testing.T) {
	s := sample.New(1)
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Fatalf("rate=1: expected Keep()=true at call %d", i)
		}
	}
	if s.Skipped() != 0 {
		t.Fatalf("expected 0 skipped, got %d", s.Skipped())
	}
}

func TestNew_ZeroRate_TreatedAsOne(t *testing.T) {
	s := sample.New(0)
	if !s.Enabled() {
		t.Fatal("rate normalised to 1 should not be Enabled")
	}
	if s.Rate() != 1 {
		t.Fatalf("expected rate 1, got %d", s.Rate())
	}
}

func TestKeep_SamplesCorrectly(t *testing.T) {
	const rate = 5
	const total = 50
	s := sample.New(rate)

	kept := 0
	for i := 0; i < total; i++ {
		if s.Keep() {
			kept++
		}
	}

	expectedKept := total / rate
	if kept != expectedKept {
		t.Fatalf("expected %d kept, got %d", expectedKept, kept)
	}
	expectedSkipped := uint64(total - expectedKept)
	if s.Skipped() != expectedSkipped {
		t.Fatalf("expected %d skipped, got %d", expectedSkipped, s.Skipped())
	}
}

func TestEnabled_ReportsCorrectly(t *testing.T) {
	if sample.New(1).Enabled() {
		t.Fatal("rate=1 should not be Enabled")
	}
	if !sample.New(2).Enabled() {
		t.Fatal("rate=2 should be Enabled")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	s := sample.New(3)
	for i := 0; i < 9; i++ {
		s.Keep()
	}
	s.Reset()
	if s.Skipped() != 0 {
		t.Fatalf("after Reset expected 0 skipped, got %d", s.Skipped())
	}
	// First call after reset must be kept (counter restarts at 1).
	if !s.Keep() {
		t.Fatal("first Keep after Reset should return true")
	}
}
