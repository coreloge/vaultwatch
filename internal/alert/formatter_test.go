package alert

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func samplePayload() Payload {
	return Payload{
		LeaseID:   "lease-abc",
		Secret:    "secret/myapp/db",
		ExpiresAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		TTL:       600,
		Severity:  SeverityWarning,
		Message:   "lease lease-abc expires in 600 seconds",
		Timestamp: time.Date(2024, 6, 1, 11, 50, 0, 0, time.UTC),
	}
}

func TestFormatter_JSONOutput(t *testing.T) {
	f := NewFormatter("json")
	data, err := f.Encode(samplePayload())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out Payload
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
	if out.LeaseID != "lease-abc" {
		t.Errorf("expected lease-abc, got %s", out.LeaseID)
	}
	if out.Severity != SeverityWarning {
		t.Errorf("expected warning severity, got %s", out.Severity)
	}
}

func TestFormatter_TextOutput(t *testing.T) {
	f := NewFormatter("text")
	data, err := f.Encode(samplePayload())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := string(data)
	if !strings.Contains(result, "[warning]") {
		t.Errorf("expected [warning] in output, got: %s", result)
	}
	if !strings.Contains(result, "lease-abc") {
		t.Errorf("expected lease-abc in output, got: %s", result)
	}
	if !strings.Contains(result, "secret/myapp/db") {
		t.Errorf("expected secret path in output, got: %s", result)
	}
}

func TestFormatter_DefaultsToJSON(t *testing.T) {
	f := NewFormatter("unknown-format")
	if f.format != FormatJSON {
		t.Errorf("expected json default, got %s", f.format)
	}
}
