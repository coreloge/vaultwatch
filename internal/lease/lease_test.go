package lease_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
)

func newInfo(ttl time.Duration) lease.Info {
	now := time.Now()
	return lease.Info{
		LeaseID:   "test/lease/123",
		Path:      "secret/data/myapp",
		Renewable: true,
		TTL:       ttl,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}
}

func TestClassify_OK(t *testing.T) {
	l := newInfo(2 * time.Hour)
	status := lease.Classify(l, 1*time.Hour, 30*time.Minute)
	if status != lease.StatusOK {
		t.Errorf("expected StatusOK, got %v", status)
	}
}

func TestClassify_Warning(t *testing.T) {
	l := newInfo(45 * time.Minute)
	status := lease.Classify(l, 1*time.Hour, 30*time.Minute)
	if status != lease.StatusWarning {
		t.Errorf("expected StatusWarning, got %v", status)
	}
}

func TestClassify_Critical(t *testing.T) {
	l := newInfo(10 * time.Minute)
	status := lease.Classify(l, 1*time.Hour, 30*time.Minute)
	if status != lease.StatusCritical {
		t.Errorf("expected StatusCritical, got %v", status)
	}
}

func TestClassify_Expired(t *testing.T) {
	l := newInfo(-1 * time.Minute)
	status := lease.Classify(l, 1*time.Hour, 30*time.Minute)
	if status != lease.StatusExpired {
		t.Errorf("expected StatusExpired, got %v", status)
	}
}

func TestRemaining_Positive(t *testing.T) {
	l := newInfo(1 * time.Hour)
	if l.Remaining() <= 0 {
		t.Error("expected positive remaining duration")
	}
}

func TestRemaining_Expired(t *testing.T) {
	l := newInfo(-5 * time.Minute)
	if l.Remaining() != 0 {
		t.Errorf("expected 0 for expired lease, got %v", l.Remaining())
	}
}

func TestIsExpired(t *testing.T) {
	expired := newInfo(-1 * time.Minute)
	if !expired.IsExpired() {
		t.Error("expected lease to be expired")
	}
	valid := newInfo(1 * time.Hour)
	if valid.IsExpired() {
		t.Error("expected lease to not be expired")
	}
}
