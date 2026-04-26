// Package backoff provides exponential backoff with jitter for retry logic.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Config holds the configuration for exponential backoff.
type Config struct {
	// InitialInterval is the starting delay before the first retry.
	InitialInterval time.Duration
	// MaxInterval is the upper bound on the delay between retries.
	MaxInterval time.Duration
	// Multiplier is applied to the interval after each attempt.
	Multiplier float64
	// MaxAttempts is the maximum number of retry attempts (0 means unlimited).
	MaxAttempts int
	// Jitter adds random noise to prevent thundering herd.
	Jitter bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		MaxAttempts:     5,
		Jitter:          true,
	}
}

// Backoff computes the delay for a given attempt number (zero-indexed).
// It applies exponential growth capped at MaxInterval, with optional jitter.
func (c Config) Backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	interval := float64(c.InitialInterval) * math.Pow(c.Multiplier, float64(attempt))
	max := float64(c.MaxInterval)
	if interval > max {
		interval = max
	}
	if c.Jitter {
		// Apply up to ±25% jitter.
		jitter := (rand.Float64()*0.5 - 0.25) * interval
		interval += jitter
		if interval < 0 {
			interval = 0
		}
	}
	return time.Duration(interval)
}

// Exceeded reports whether the attempt count has surpassed MaxAttempts.
// If MaxAttempts is 0, it never exceeds.
func (c Config) Exceeded(attempt int) bool {
	if c.MaxAttempts == 0 {
		return false
	}
	return attempt >= c.MaxAttempts
}
