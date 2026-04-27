// Package sampler provides probabilistic sampling for lease alert events.
// It allows reducing alert volume by only forwarding a configurable
// fraction of events that match a given status classification.
package sampler

import (
	"math/rand"
	"sync"

	"github.com/yourusername/vaultwatch/internal/lease"
)

// Config holds sampler configuration.
type Config struct {
	// Rate is the fraction of events to allow through, in the range [0.0, 1.0].
	// A rate of 1.0 allows all events; 0.0 drops all events.
	Rate float64

	// Statuses restricts sampling to specific lease statuses.
	// If empty, all statuses are subject to sampling.
	Statuses []lease.Status
}

// DefaultConfig returns a Config that passes all events through.
func DefaultConfig() Config {
	return Config{Rate: 1.0}
}

// Sampler probabilistically allows or drops lease events.
type Sampler struct {
	mu       sync.Mutex
	rate     float64
	statuses map[lease.Status]struct{}
	rng      *rand.Rand
}

// New creates a new Sampler from the given Config.
// If cfg.Rate is outside [0.0, 1.0] it is clamped.
func New(cfg Config, src rand.Source) *Sampler {
	rate := cfg.Rate
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}

	statuses := make(map[lease.Status]struct{}, len(cfg.Statuses))
	for _, s := range cfg.Statuses {
		statuses[s] = struct{}{}
	}

	if src == nil {
		src = rand.NewSource(42)
	}

	return &Sampler{
		rate:     rate,
		statuses: statuses,
		rng:      rand.New(src),
	}
}

// Allow returns true if the event should be forwarded.
// Events whose status is not in the configured set are always allowed.
func (s *Sampler) Allow(info lease.Info) bool {
	if len(s.statuses) > 0 {
		if _, ok := s.statuses[info.Status]; !ok {
			return true
		}
	}

	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()

	return v < s.rate
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
