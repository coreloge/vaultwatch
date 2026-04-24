package alert

import (
	"fmt"
	"time"
)

// Severity represents the urgency level of an alert.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Payload is the structured alert sent to webhook endpoints.
type Payload struct {
	LeaseID   string    `json:"lease_id"`
	Secret    string    `json:"secret"`
	ExpiresAt time.Time `json:"expires_at"`
	TTL       int64     `json:"ttl_seconds"`
	Severity  Severity  `json:"severity"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Builder constructs alert payloads based on lease TTL thresholds.
type Builder struct {
	criticalThreshold time.Duration
	warningThreshold  time.Duration
}

// NewBuilder creates a Builder with the given warning and critical thresholds.
func NewBuilder(warningThreshold, criticalThreshold time.Duration) *Builder {
	return &Builder{
		warningThreshold:  warningThreshold,
		criticalThreshold: criticalThreshold,
	}
}

// Build creates a Payload for the given lease details.
func (b *Builder) Build(leaseID, secret string, expiresAt time.Time, ttl int64) Payload {
	remaining := time.Until(expiresAt)
	severity := b.classify(remaining)

	return Payload{
		LeaseID:   leaseID,
		Secret:    secret,
		ExpiresAt: expiresAt,
		TTL:       ttl,
		Severity:  severity,
		Message:   fmt.Sprintf("lease %s expires in %.0f seconds", leaseID, remaining.Seconds()),
		Timestamp: time.Now().UTC(),
	}
}

func (b *Builder) classify(remaining time.Duration) Severity {
	if remaining <= b.criticalThreshold {
		return SeverityCritical
	}
	if remaining <= b.warningThreshold {
		return SeverityWarning
	}
	return SeverityInfo
}
