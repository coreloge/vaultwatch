package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/audit"
	"github.com/yourusername/vaultwatch/internal/lease"
)

func newMiddleware() (*audit.LeaseEventLogger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	l := audit.New(buf)
	return audit.NewLeaseEventLogger(l), buf
}

func sampleInfo() lease.Info {
	return lease.Info{
		LeaseID:   "secret/data/db#abc123",
		Path:      "secret/data/db",
		TTL:       300,
		ExpireTime: time.Now().Add(5 * time.Minute),
	}
}

func TestOnLeaseChecked_WritesEvent(t *testing.T) {
	m, buf := newMiddleware()

	if err := m.OnLeaseChecked(sampleInfo()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Type != audit.EventLeaseChecked {
		t.Errorf("expected %q, got %q", audit.EventLeaseChecked, e.Type)
	}
	if e.Meta["path"] != "secret/data/db" {
		t.Errorf("expected path in meta, got %v", e.Meta)
	}
}

func TestOnAlertSent_WritesEvent(t *testing.T) {
	m, buf := newMiddleware()

	_ = m.OnAlertSent("lease-1", "http://hook.example.com")

	var e audit.Event
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Type != audit.EventAlertSent {
		t.Errorf("expected %q, got %q", audit.EventAlertSent, e.Type)
	}
	if e.Meta["webhook"] != "http://hook.example.com" {
		t.Errorf("expected webhook in meta, got %v", e.Meta)
	}
}

func TestOnAlertFailed_WritesEvent(t *testing.T) {
	m, buf := newMiddleware()

	_ = m.OnAlertFailed("lease-2", "connection refused")

	var e audit.Event
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Type != audit.EventAlertFailed {
		t.Errorf("expected %q, got %q", audit.EventAlertFailed, e.Type)
	}
	if e.Meta["reason"] != "connection refused" {
		t.Errorf("expected reason in meta, got %v", e.Meta)
	}
}

func TestOnLeaseRenewed_WritesEvent(t *testing.T) {
	m, buf := newMiddleware()

	_ = m.OnLeaseRenewed("lease-3")

	var e audit.Event
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Type != audit.EventLeaseRenewed {
		t.Errorf("expected %q, got %q", audit.EventLeaseRenewed, e.Type)
	}
	if e.LeaseID != "lease-3" {
		t.Errorf("expected lease_id lease-3, got %q", e.LeaseID)
	}
}
