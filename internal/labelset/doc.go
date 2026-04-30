// Package labelset provides a lightweight immutable key-value label container
// used to attach structured metadata to lease events, alert payloads, and
// audit log entries throughout vaultwatch.
//
// Labels flow from lease discovery through the processing pipeline and are
// attached to outbound webhook payloads so downstream consumers can filter
// and route alerts without parsing raw message bodies.
package labelset
