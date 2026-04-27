package redact_test

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/redact"
)

func newRedactor(extra ...string) *redact.Redactor {
	return redact.New(extra...)
}

func TestMap_RedactsSensitiveKeys(t *testing.T) {
	r := newRedactor()
	input := map[string]string{
		"vault_token": "s.supersecret",
		"db_password": "hunter2",
		"lease_id":    "lease/abc/123",
	}
	out := r.Map(input)

	if out["vault_token"] != "[REDACTED]" {
		t.Errorf("expected vault_token to be redacted, got %q", out["vault_token"])
	}
	if out["db_password"] != "[REDACTED]" {
		t.Errorf("expected db_password to be redacted, got %q", out["db_password"])
	}
	if out["lease_id"] != "lease/abc/123" {
		t.Errorf("expected lease_id to be unchanged, got %q", out["lease_id"])
	}
}

func TestMap_PreservesNonSensitiveKeys(t *testing.T) {
	r := newRedactor()
	input := map[string]string{
		"mount":    "secret/",
		"duration": "72h",
	}
	out := r.Map(input)

	for k, want := range input {
		if out[k] != want {
			t.Errorf("key %q: got %q, want %q", k, out[k], want)
		}
	}
}

func TestMap_DoesNotMutateOriginal(t *testing.T) {
	r := newRedactor()
	input := map[string]string{"api_key": "abc123"}
	_ = r.Map(input)
	if input["api_key"] != "abc123" {
		t.Error("original map was mutated")
	}
}

func TestValue_SensitiveKey(t *testing.T) {
	r := newRedactor()
	got := r.Value("auth_header", "Bearer tok")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestValue_NonSensitiveKey(t *testing.T) {
	r := newRedactor()
	got := r.Value("lease_id", "lease/abc")
	if got != "lease/abc" {
		t.Errorf("expected lease/abc, got %q", got)
	}
}

func TestNew_ExtraKeys(t *testing.T) {
	r := newRedactor("fingerprint")
	got := r.Value("device_fingerprint", "abc")
	if got != "[REDACTED]" {
		t.Errorf("expected extra key to be redacted, got %q", got)
	}
}

func TestMap_CaseInsensitiveMatch(t *testing.T) {
	r := newRedactor()
	input := map[string]string{"API_SECRET": "topsecret"}
	out := r.Map(input)
	if out["API_SECRET"] != "[REDACTED]" {
		t.Errorf("expected case-insensitive redaction, got %q", out["API_SECRET"])
	}
}
