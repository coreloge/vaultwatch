package metrics_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/metrics"
)

func newTestMetrics(t *testing.T) *metrics.Metrics {
	t.Helper()
	m := metrics.New()
	if m == nil {
		t.Fatal("expected non-nil Metrics")
	}
	return m
}

func TestNew_ReturnsNonNil(t *testing.T) {
	m := metrics.New()
	if m == nil {
		t.Error("New() returned nil")
	}
}

func TestRecordLeaseChecked(t *testing.T) {
	m := newTestMetrics(t)

	m.RecordLeaseChecked()
	m.RecordLeaseChecked()
	m.RecordLeaseChecked()

	snap := m.Snapshot()
	if snap.LeasesChecked != 3 {
		t.Errorf("expected LeasesChecked=3, got %d", snap.LeasesChecked)
	}
}

func TestRecordAlertSent(t *testing.T) {
	m := newTestMetrics(t)

	m.RecordAlertSent("warning")
	m.RecordAlertSent("critical")
	m.RecordAlertSent("warning")

	snap := m.Snapshot()
	if snap.AlertsSent != 3 {
		t.Errorf("expected AlertsSent=3, got %d", snap.AlertsSent)
	}
	if snap.AlertsByLevel["warning"] != 2 {
		t.Errorf("expected 2 warning alerts, got %d", snap.AlertsByLevel["warning"])
	}
	if snap.AlertsByLevel["critical"] != 1 {
		t.Errorf("expected 1 critical alert, got %d", snap.AlertsByLevel["critical"])
	}
}

func TestRecordWebhookError(t *testing.T) {
	m := newTestMetrics(t)

	m.RecordWebhookError()
	m.RecordWebhookError()

	snap := m.Snapshot()
	if snap.WebhookErrors != 2 {
		t.Errorf("expected WebhookErrors=2, got %d", snap.WebhookErrors)
	}
}

func TestRecordLeaseExpired(t *testing.T) {
	m := newTestMetrics(t)

	m.RecordLeaseExpired()

	snap := m.Snapshot()
	if snap.LeasesExpired != 1 {
		t.Errorf("expected LeasesExpired=1, got %d", snap.LeasesExpired)
	}
}

func TestSetLastCheckTime(t *testing.T) {
	m := newTestMetrics(t)

	now := time.Now().UTC().Truncate(time.Second)
	m.SetLastCheckTime(now)

	snap := m.Snapshot()
	if !snap.LastCheckTime.Equal(now) {
		t.Errorf("expected LastCheckTime=%v, got %v", now, snap.LastCheckTime)
	}
}

func TestSnapshot_IsImmutable(t *testing.T) {
	m := newTestMetrics(t)

	m.RecordLeaseChecked()
	snap1 := m.Snapshot()

	m.RecordLeaseChecked()
	snap2 := m.Snapshot()

	if snap1.LeasesChecked == snap2.LeasesChecked {
		t.Error("expected snapshots to reflect independent state")
	}
	if snap1.LeasesChecked != 1 {
		t.Errorf("snap1 should still show 1, got %d", snap1.LeasesChecked)
	}
	if snap2.LeasesChecked != 2 {
		t.Errorf("snap2 should show 2, got %d", snap2.LeasesChecked)
	}
}
