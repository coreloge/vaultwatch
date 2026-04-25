// Package audit provides structured, newline-delimited JSON audit logging
// for VaultWatch lease and alert lifecycle events.
//
// Usage:
//
//	logger := audit.New(os.Stderr)
//	el := audit.NewLeaseEventLogger(logger)
//	el.OnLeaseChecked(info)
//	el.OnAlertSent(leaseID, webhookURL)
//
// Each event is written as a single JSON object followed by a newline,
// suitable for ingestion by log aggregators such as Loki or Splunk.
package audit
