package jitter_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/jitter"
)

// fixedSource always returns the same value so tests are deterministic.
type fixedSource struct{ v float64 }

func (f *fixedSource) Float64() float64 { return f.v }

func newJitterer(factor, srcVal float64) *jitter.Jitterer {
	return jitter.New(factor, &fixedSource{v: srcVal})
}

func TestApply_AddsPositiveJitter(t *testing.T) {
	j := newJitterer(0.5, 1.0) // src always returns 1.0 → max jitter
	base := 10 * time.Second
	got := j.Apply(base)
	// expect base + 0.5*base*1.0 = 15s
	want := 15 * time.Second
	if got != want {
		t.Errorf("Apply() = %v, want %v", got, want)
	}
}

func TestApply_ZeroSrcReturnsBase(t *testing.T) {
	j := newJitterer(0.5, 0.0) // src returns 0 → no jitter added
	base := 10 * time.Second
	got := j.Apply(base)
	if got != base {
		t.Errorf("Apply() = %v, want %v", got, base)
	}
}

func TestApply_NegativeBasePassedThrough(t *testing.T) {
	j := newJitterer(0.25, 0.5)
	got := j.Apply(-1 * time.Second)
	if got != -1*time.Second {
		t.Errorf("Apply() = %v, want unchanged negative base", got)
	}
}

func TestApplySigned_CanReduceBase(t *testing.T) {
	// src=0.0 → delta = (0*2-1)*half = -half → result < base
	j := newJitterer(0.5, 0.0)
	base := 10 * time.Second
	got := j.ApplySigned(base)
	if got >= base {
		t.Errorf("ApplySigned() = %v, expected less than base %v", got, base)
	}
}

func TestApplySigned_CanIncreaseBase(t *testing.T) {
	// src=1.0 → delta = (2-1)*half = +half → result > base
	j := newJitterer(0.5, 1.0)
	base := 10 * time.Second
	got := j.ApplySigned(base)
	if got <= base {
		t.Errorf("ApplySigned() = %v, expected greater than base %v", got, base)
	}
}

func TestNew_ClampsFactor(t *testing.T) {
	// factor > 1 should be clamped to 1; result must not exceed 2*base
	j := newJitterer(99, 1.0)
	base := 10 * time.Second
	got := j.Apply(base)
	if got > 2*base {
		t.Errorf("Apply() = %v exceeds 2*base after factor clamp", got)
	}
}

func TestNew_NilSourceUsesDefault(t *testing.T) {
	// Ensure no panic when src is nil (uses internal default).
	j := jitter.New(0.1, nil)
	base := 5 * time.Second
	got := j.Apply(base)
	if got < base {
		t.Errorf("Apply() = %v, should be >= base %v", got, base)
	}
}
