// Package envelope provides the Envelope type which wraps a lease.Info
// value with delivery metadata including a unique message ID, creation
// timestamp, attempt counter, and origin label.
//
// Envelopes are created once per alert event and passed through the
// delivery pipeline so that retry logic, audit logging, and deduplication
// components can share a stable identifier without re-deriving it.
package envelope
