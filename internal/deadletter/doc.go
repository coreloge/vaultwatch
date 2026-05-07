// Package deadletter implements a bounded dead-letter store for VaultWatch.
//
// When a webhook alert cannot be delivered after all retry attempts are
// exhausted the calling component should record the failure here so that
// operators can inspect, replay, or export the undelivered events.
//
// Entries are automatically evicted once their TTL elapses, preventing
// unbounded memory growth in long-running deployments.
package deadletter
