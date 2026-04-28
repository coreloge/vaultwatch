// Package escalation implements multi-tier alert escalation for VaultWatch.
//
// When a lease remains in a degraded state (warning or critical) for an
// extended period, the escalation tier increases from TierWarning through
// TierCritical to TierEmergency. Callers can use the tier value to route
// alerts to increasingly urgent notification channels.
//
// Usage:
//
//	esc := escalation.New(escalation.DefaultConfig())
//	tier := esc.Evaluate(leaseInfo)
//	if tier >= escalation.TierCritical {
//		// page on-call
//	}
package escalation
