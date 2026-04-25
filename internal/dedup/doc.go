// Package dedup implements alert deduplication for VaultWatch.
//
// A Deduplicator tracks the most-recently alerted (leaseID, status) pair
// and suppresses repeat notifications that arrive within a configurable
// time window, reducing noise when a lease stays in the same degraded
// state across multiple monitor ticks.
//
// Usage:
//
//	dd := dedup.New(5 * time.Minute)
//	if !dd.IsDuplicate(info.LeaseID, info.Status) {
//		// send alert
//		dd.Record(info.LeaseID, info.Status)
//	}
package dedup
