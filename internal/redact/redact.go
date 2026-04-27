// Package redact provides utilities for scrubbing sensitive values
// from lease metadata and alert payloads before they are logged or
// transmitted via webhook.
package redact

import (
	"strings"
)

const masked = "[REDACTED]"

// sensitiveKeys contains substrings that indicate a field holds a
// secret value and should be masked before leaving the process.
var sensitiveKeys = []string{
	"token",
	"password",
	"secret",
	"key",
	"credential",
	"auth",
	"passwd",
}

// Redactor scrubs sensitive fields from arbitrary string maps.
type Redactor struct {
	extraKeys []string
}

// New returns a Redactor. Additional key substrings to treat as
// sensitive may be supplied via extraKeys.
func New(extraKeys ...string) *Redactor {
	return &Redactor{extraKeys: extraKeys}
}

// Map returns a shallow copy of m with sensitive values replaced by
// the masked placeholder. Keys are matched case-insensitively.
func (r *Redactor) Map(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if r.isSensitive(k) {
			out[k] = masked
		} else {
			out[k] = v
		}
	}
	return out
}

// Value returns the masked placeholder when key is considered
// sensitive, otherwise it returns value unchanged.
func (r *Redactor) Value(key, value string) string {
	if r.isSensitive(key) {
		return masked
	}
	return value
}

// isSensitive reports whether key contains any known sensitive
// substring (case-insensitive).
func (r *Redactor) isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	for _, s := range r.extraKeys {
		if strings.Contains(lower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}
