// Package lease provides types and utilities for managing Vault secret leases.
package lease

import (
	"time"
)

// ExpiryWindow defines the thresholds used to classify lease urgency.
type ExpiryWindow struct {
	Critical time.Duration
	Warning  time.Duration
}

// DefaultExpiryWindow returns the standard expiry classification thresholds.
func DefaultExpiryWindow() ExpiryWindow {
	return ExpiryWindow{
		Critical: 1 * time.Hour,
		Warning:  6 * time.Hour,
	}
}

// ExpiryChecker evaluates lease TTLs against configurable thresholds.
type ExpiryChecker struct {
	window ExpiryWindow
	now    func() time.Time
}

// NewExpiryChecker creates an ExpiryChecker with the given window.
func NewExpiryChecker(w ExpiryWindow) *ExpiryChecker {
	return &ExpiryChecker{
		window: w,
		now:    time.Now,
	}
}

// IsCritical reports whether the TTL falls within the critical threshold.
func (c *ExpiryChecker) IsCritical(ttl TTL) bool {
	remaining := ttl.RemainingFrom(c.now())
	return remaining >= 0 && remaining <= c.window.Critical
}

// IsWarning reports whether the TTL falls within the warning threshold
// but outside the critical threshold.
func (c *ExpiryChecker) IsWarning(ttl TTL) bool {
	remaining := ttl.RemainingFrom(c.now())
	return remaining > c.window.Critical && remaining <= c.window.Warning
}

// IsExpired reports whether the TTL has already elapsed.
func (c *ExpiryChecker) IsExpired(ttl TTL) bool {
	return ttl.RemainingFrom(c.now()) < 0
}

// StatusFor returns the Status classification for the given TTL.
func (c *ExpiryChecker) StatusFor(ttl TTL) Status {
	switch {
	case c.IsExpired(ttl):
		return StatusExpired
	case c.IsCritical(ttl):
		return StatusCritical
	case c.IsWarning(ttl):
		return StatusWarning
	default:
		return StatusOK
	}
}
