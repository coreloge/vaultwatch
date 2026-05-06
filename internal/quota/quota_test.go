package quota_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/quota"
)

func newQuota(max int, window time.Duration) *quota.Quota {
	return quota.New(quota.Config{MaxAlerts: max, Window: window})
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	q := newQuota(3, time.Minute)
	if !q.Allow("secret/prod/") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_BlockedAfterMax(t *testing.T) {
	q := newQuota(2, time.Minute)
	if !q.Allow("secret/prod/") {
		t.Fatal("expected first call allowed")
	}
	if !q.Allow("secret/prod/") {
		t.Fatal("expected second call allowed")
	}
	if q.Allow("secret/prod/") {
		t.Fatal("expected third call to be blocked")
	}
}

func TestAllow_DifferentPrefixesIndependent(t *testing.T) {
	q := newQuota(1, time.Minute)
	q.Allow("secret/prod/")
	if !q.Allow("secret/staging/") {
		t.Fatal("expected different prefix to be allowed independently")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	q := newQuota(1, 10*time.Millisecond)
	q.Allow("secret/prod/")
	if q.Allow("secret/prod/") {
		t.Fatal("expected second call within window to be blocked")
	}
	time.Sleep(20 * time.Millisecond)
	if !q.Allow("secret/prod/") {
		t.Fatal("expected call after window reset to be allowed")
	}
}

func TestRemaining_DecreasesWithEachAllow(t *testing.T) {
	q := newQuota(3, time.Minute)
	remaining, _ := q.Remaining("secret/prod/")
	if remaining != 3 {
		t.Fatalf("expected 3 remaining before any allow, got %d", remaining)
	}
	q.Allow("secret/prod/")
	remaining, _ = q.Remaining("secret/prod/")
	if remaining != 2 {
		t.Fatalf("expected 2 remaining after one allow, got %d", remaining)
	}
}

func TestRemaining_NeverNegative(t *testing.T) {
	q := newQuota(1, time.Minute)
	q.Allow("secret/prod/")
	q.Allow("secret/prod/") // blocked but count should not go negative
	remaining, _ := q.Remaining("secret/prod/")
	if remaining < 0 {
		t.Fatalf("remaining should not be negative, got %d", remaining)
	}
}

func TestReset_ClearsQuota(t *testing.T) {
	q := newQuota(1, time.Minute)
	q.Allow("secret/prod/")
	if q.Allow("secret/prod/") {
		t.Fatal("expected second call blocked before reset")
	}
	q.Reset("secret/prod/")
	if !q.Allow("secret/prod/") {
		t.Fatal("expected call to be allowed after reset")
	}
}

func TestDefaultConfig_SaneValues(t *testing.T) {
	cfg := quota.DefaultConfig()
	if cfg.MaxAlerts <= 0 {
		t.Errorf("expected positive MaxAlerts, got %d", cfg.MaxAlerts)
	}
	if cfg.Window <= 0 {
		t.Errorf("expected positive Window, got %v", cfg.Window)
	}
}
