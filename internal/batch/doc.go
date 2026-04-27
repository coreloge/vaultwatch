// Package batch provides a time-and-size-bounded collector for lease events.
//
// A Collector accumulates lease.Info values and flushes them as a grouped
// slice to a registered Handler function. Flushing is triggered either when
// the number of buffered items reaches MaxSize or when the configured Window
// duration elapses since the first item was added.
//
// This is useful for downstream systems that prefer bulk ingestion over
// per-event delivery, reducing webhook round-trips during bursts of
// expiring leases.
package batch
