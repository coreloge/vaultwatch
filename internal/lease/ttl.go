// Package lease provides types and utilities for managing Vault secret leases.
package lease

import (
	"fmt"
	"time"
)

// TTL represents a lease time-to-live duration with helper methods.
type TTL struct {
	value time.Duration
}

// NewTTL creates a TTL from a duration.
func NewTTL(d time.Duration) TTL {
	return TTL{value: d}
}

// NewTTLFromSeconds creates a TTL from a seconds value as returned by Vault.
func NewTTLFromSeconds(seconds int64) TTL {
	return TTL{value: time.Duration(seconds) * time.Second}
}

// Duration returns the underlying time.Duration.
func (t TTL) Duration() time.Duration {
	return t.value
}

// Seconds returns the TTL in whole seconds.
func (t TTL) Seconds() int64 {
	return int64(t.value.Seconds())
}

// IsZero reports whether the TTL is zero or negative.
func (t TTL) IsZero() bool {
	return t.value <= 0
}

// ExpiresAt returns the absolute expiry time given a reference time.
func (t TTL) ExpiresAt(from time.Time) time.Time {
	return from.Add(t.value)
}

// String returns a human-readable representation of the TTL.
func (t TTL) String() string {
	if t.IsZero() {
		return "expired"
	}
	h := int(t.value.Hours())
	m := int(t.value.Minutes()) % 60
	s := int(t.value.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// RemainingFrom returns the remaining TTL relative to now, clamped to zero.
func (t TTL) RemainingFrom(reference time.Time) TTL {
	remaining := time.Until(t.ExpiresAt(reference))
	if remaining < 0 {
		remaining = 0
	}
	return TTL{value: remaining}
}
