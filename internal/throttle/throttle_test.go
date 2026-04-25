package throttle_test

import (
	"testing"
	"time"

	"github.com/vaultwatch/internal/throttle"
)

func newThrottler(window time.Duration) *throttle.Throttler {
	return throttle.New(window)
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	th := newThrottler(time.Minute)
	if !th.Allow("lease-1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	th := newThrottler(time.Minute)
	th.Allow("lease-1")
	if th.Allow("lease-1") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllow_DifferentLeasesIndependent(t *testing.T) {
	th := newThrottler(time.Minute)
	th.Allow("lease-1")
	if !th.Allow("lease-2") {
		t.Fatal("expected different lease ID to be allowed independently")
	}
}

func TestAllow_PermittedAfterWindow(t *testing.T) {
	th := newThrottler(10 * time.Millisecond)
	th.Allow("lease-1")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("lease-1") {
		t.Fatal("expected call to be allowed after window elapsed")
	}
}

func TestReset_ClearsState(t *testing.T) {
	th := newThrottler(time.Minute)
	th.Allow("lease-1")
	th.Reset("lease-1")
	if !th.Allow("lease-1") {
		t.Fatal("expected allow after reset")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	th := newThrottler(10 * time.Millisecond)
	th.Allow("lease-1")
	th.Allow("lease-2")
	time.Sleep(20 * time.Millisecond)
	th.Purge()
	if th.Len() != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", th.Len())
	}
}

func TestPurge_RetainsActiveEntries(t *testing.T) {
	th := newThrottler(time.Minute)
	th.Allow("lease-1")
	th.Purge()
	if th.Len() != 1 {
		t.Fatalf("expected 1 active entry retained, got %d", th.Len())
	}
}

func TestLen_ReflectsTrackedLeases(t *testing.T) {
	th := newThrottler(time.Minute)
	if th.Len() != 0 {
		t.Fatal("expected empty throttler to have len 0")
	}
	th.Allow("a")
	th.Allow("b")
	if th.Len() != 2 {
		t.Fatalf("expected len 2, got %d", th.Len())
	}
}
