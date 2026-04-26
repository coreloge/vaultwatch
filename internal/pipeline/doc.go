// Package pipeline provides a composable lease-event processing pipeline
// that chains together deduplication, suppression, filtering, throttling,
// and alert dispatch into a single reusable unit.
//
// A Pipeline is constructed with a set of middleware-style components and
// a terminal Dispatcher. When Process is called with a LeaseInfo, the event
// travels through each stage in order:
//
//  1. Deduplication  – drops events whose status has not changed since the
//     last observation within the configured window.
//  2. Suppression    – drops events for leases that have been explicitly
//     silenced by an operator.
//  3. Filtering      – drops events that do not match the include/exclude
//     rules defined in configuration.
//  4. Throttling     – rate-limits repeated alerts for the same lease so
//     that noisy leases do not flood downstream webhooks.
//  5. Dispatch       – formats and delivers the alert payload to every
//     configured webhook endpoint.
//
// Each stage is optional; passing a nil component causes that stage to be
// skipped transparently.
package pipeline
