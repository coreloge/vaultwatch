// Package rollup provides time-window batching for lease alert events.
//
// Instead of dispatching a webhook for every individual lease event,
// rollup collects events over a configurable window (or up to a maximum
// batch size) and delivers them as a single Batch. This reduces downstream
// webhook pressure when many leases expire or change status simultaneously.
//
// Usage:
//
//	r := rollup.New(rollup.DefaultConfig())
//	go func() {
//		for batch := range r.Batches() {
//			// handle batch.Events
//		}
//	}()
//	r.Add(leaseInfo)
package rollup
