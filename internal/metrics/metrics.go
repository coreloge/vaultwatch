// Package metrics provides lightweight in-process counters and gauges
// for tracking VaultWatch operational statistics. These can be scraped
// or logged periodically without requiring an external metrics backend.
package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics holds all runtime counters and gauges for the daemon.
type Metrics struct {
	mu sync.RWMutex

	// Lease counters
	LeasesChecked   atomic.Int64
	LeasesExpired   atomic.Int64
	LeasesWarning   atomic.Int64
	LeasesCritical  atomic.Int64

	// Webhook counters
	WebhooksSent    atomic.Int64
	WebhooksFailed  atomic.Int64
	WebhooksRetried atomic.Int64

	// Timing
	LastCheckAt  time.Time
	StartedAt    time.Time
}

// New creates and initialises a Metrics instance, recording the current
// time as the daemon start time.
func New() *Metrics {
	return &Metrics{
		StartedAt: time.Now(),
	}
}

// RecordCheck increments the total lease check counter and updates the
// timestamp of the most recent check cycle.
func (m *Metrics) RecordCheck() {
	m.LeasesChecked.Add(1)
	m.mu.Lock()
	m.LastCheckAt = time.Now()
	m.mu.Unlock()
}

// RecordLeaseStatus increments the appropriate severity counter based on
// the string label returned by the lease classifier ("warning", "critical",
// "expired"; anything else is a no-op).
func (m *Metrics) RecordLeaseStatus(status string) {
	switch status {
	case "warning":
		m.LeasesWarning.Add(1)
	case "critical":
		m.LeasesCritical.Add(1)
	case "expired":
		m.LeasesExpired.Add(1)
	}
}

// RecordWebhookSent increments the successful webhook delivery counter.
func (m *Metrics) RecordWebhookSent() {
	m.WebhooksSent.Add(1)
}

// RecordWebhookFailed increments the failed webhook delivery counter.
func (m *Metrics) RecordWebhookFailed() {
	m.WebhooksFailed.Add(1)
}

// RecordWebhookRetry increments the webhook retry counter.
func (m *Metrics) RecordWebhookRetry() {
	m.WebhooksRetried.Add(1)
}

// Snapshot returns a point-in-time copy of all metric values, safe for
// serialisation or logging without holding any locks.
func (m *Metrics) Snapshot() Snapshot {
	m.mu.RLock()
	lastCheck := m.LastCheckAt
	m.mu.RUnlock()

	return Snapshot{
		LeasesChecked:   m.LeasesChecked.Load(),
		LeasesExpired:   m.LeasesExpired.Load(),
		LeasesWarning:   m.LeasesWarning.Load(),
		LeasesCritical:  m.LeasesCritical.Load(),
		WebhooksSent:    m.WebhooksSent.Load(),
		WebhooksFailed:  m.WebhooksFailed.Load(),
		WebhooksRetried: m.WebhooksRetried.Load(),
		UptimeSeconds:   int64(time.Since(m.StartedAt).Seconds()),
		LastCheckAt:     lastCheck,
	}
}

// Snapshot is an immutable point-in-time copy of Metrics values.
type Snapshot struct {
	LeasesChecked   int64     `json:"leases_checked"`
	LeasesExpired   int64     `json:"leases_expired"`
	LeasesWarning   int64     `json:"leases_warning"`
	LeasesCritical  int64     `json:"leases_critical"`
	WebhooksSent    int64     `json:"webhooks_sent"`
	WebhooksFailed  int64     `json:"webhooks_failed"`
	WebhooksRetried int64     `json:"webhooks_retried"`
	UptimeSeconds   int64     `json:"uptime_seconds"`
	LastCheckAt     time.Time `json:"last_check_at"`
}
