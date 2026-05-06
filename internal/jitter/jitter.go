// Package jitter provides utilities for adding randomised jitter to
// durations, preventing thundering-herd problems when many leases
// are checked or alerts are dispatched at the same instant.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is the interface satisfied by a random-number source.
type Source interface {
	Float64() float64
}

// defaultSource wraps the global rand functions behind a mutex so it
// is safe for concurrent use.
type defaultSource struct{ mu sync.Mutex }

func (d *defaultSource) Float64() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return rand.Float64() //nolint:gosec // non-cryptographic jitter
}

// Jitterer applies a configurable jitter factor to durations.
type Jitterer struct {
	factor float64 // fraction of the base duration, e.g. 0.25 → ±25 %
	src    Source
}

// New returns a Jitterer that adds up to factor*base random jitter.
// factor must be in the range (0, 1]; values outside that range are
// clamped silently.
func New(factor float64, src Source) *Jitterer {
	if factor <= 0 {
		factor = 0.1
	}
	if factor > 1 {
		factor = 1
	}
	if src == nil {
		src = &defaultSource{}
	}
	return &Jitterer{factor: factor, src: src}
}

// Apply returns base plus a random duration in [0, factor*base).
func (j *Jitterer) Apply(base time.Duration) time.Duration {
	if base <= 0 {
		return base
	}
	max := float64(base) * j.factor
	delta := time.Duration(j.src.Float64() * max)
	return base + delta
}

// ApplySigned returns base plus a random duration in
// [-factor*base/2, +factor*base/2), centred around base.
func (j *Jitterer) ApplySigned(base time.Duration) time.Duration {
	if base <= 0 {
		return base
	}
	half := float64(base) * j.factor / 2
	delta := time.Duration((j.src.Float64()*2 - 1) * half)
	return base + delta
}
