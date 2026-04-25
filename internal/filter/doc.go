// Package filter provides include/exclude filtering for Vault lease events
// in vaultwatch.
//
// Filters are evaluated in order: exclude rules are checked first, and if a
// lease matches any exclude rule it is rejected regardless of include rules.
// If no include rules are configured, all non-excluded leases are permitted.
//
// Rules can match on lease path prefix, status, or both. When both fields are
// set on a single Rule, both conditions must be satisfied for the rule to match.
package filter
