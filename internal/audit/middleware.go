package audit

import (
	"github.com/yourusername/vaultwatch/internal/lease"
)

// LeaseEventLogger wraps an audit Logger and exposes domain-specific helpers
// for recording lease lifecycle events.
type LeaseEventLogger struct {
	logger *Logger
}

// NewLeaseEventLogger returns a LeaseEventLogger backed by the given Logger.
func NewLeaseEventLogger(l *Logger) *LeaseEventLogger {
	return &LeaseEventLogger{logger: l}
}

// OnLeaseChecked records a lease check event with its classification.
func (m *LeaseEventLogger) OnLeaseChecked(info lease.Info) error {
	return m.logger.Log(EventLeaseChecked, info.LeaseID, map[string]string{
		"path":   info.Path,
		"status": string(lease.Classify(info)),
	})
}

// OnAlertSent records a successful alert dispatch.
func (m *LeaseEventLogger) OnAlertSent(leaseID, webhookURL string) error {
	return m.logger.Log(EventAlertSent, leaseID, map[string]string{
		"webhook": webhookURL,
	})
}

// OnAlertFailed records a failed alert dispatch.
func (m *LeaseEventLogger) OnAlertFailed(leaseID, reason string) error {
	return m.logger.Log(EventAlertFailed, leaseID, map[string]string{
		"reason": reason,
	})
}

// OnLeaseRenewed records a lease renewal event.
func (m *LeaseEventLogger) OnLeaseRenewed(leaseID string) error {
	return m.logger.Log(EventLeaseRenewed, leaseID, nil)
}
