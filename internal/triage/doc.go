// Package triage implements a priority queue for lease alert events.
//
// Leases are ranked by severity (critical > warning > info) and then
// by the time they were enqueued, ensuring the most urgent and oldest
// alerts are dispatched before lower-priority, newer ones.
//
// Typical usage:
//
//	q := triage.New()
//	q.Add(info)
//	for _, entry := range q.Drain() {
//		dispatcher.Dispatch(ctx, entry.Info)
//	}
package triage
