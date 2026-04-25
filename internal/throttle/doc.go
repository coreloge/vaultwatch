// Package throttle implements per-lease alert throttling for vaultwatch.
//
// When many leases approach expiration simultaneously, a naive alerting
// system would emit a flood of webhook calls. The Throttler enforces a
// configurable minimum window between successive alerts for the same
// lease ID, preventing alert storms while still ensuring timely
// notification for genuinely new state changes.
//
// Usage:
//
//	th := throttle.New(5 * time.Minute)
//	if th.Allow(leaseID) {
//		// send alert
//	}
package throttle
