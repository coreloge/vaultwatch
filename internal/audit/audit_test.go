package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/audit"
)

func newTestLogger() (*audit.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return audit.New(buf), buf
}

func TestLog_WritesJSON(t *testing.T) {
	logger, buf := newTestLogger()

	err := logger.Log(audit.EventLeaseChecked, "lease-123", map[string]string{"status": "warning"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var event audit.Event
	if err := json.Unmarshal(buf.Bytes(), &event); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if event.Type != audit.EventLeaseChecked {
		t.Errorf("expected type %q, got %q", audit.EventLeaseChecked, event.Type)
	}
	if event.LeaseID != "lease-123" {
		t.Errorf("expected lease_id %q, got %q", "lease-123", event.LeaseID)
	}
	if event.Meta["status"] != "warning" {
		t.Errorf("expected meta status=warning, got %v", event.Meta)
	}
}

func TestLog_NewlineDelimited(t *testing.T) {
	logger, buf := newTestLogger()

	_ = logger.Log(audit.EventAlertSent, "a", nil)
	_ = logger.Log(audit.EventAlertFailed, "b", nil)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestLog_TimestampPresent(t *testing.T) {
	logger, buf := newTestLogger()
	_ = logger.Log(audit.EventLeaseRenewed, "x", nil)

	var event audit.Event
	_ = json.Unmarshal(buf.Bytes(), &event)
	if event.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNew_NilWriterUsesStdout(t *testing.T) {
	logger := audit.New(nil)
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}
