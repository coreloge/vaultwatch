package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/ratelimit"
)

func newLimiter(d time.Duration) *ratelimit.Limiter {
	return ratelimit.New(d)
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := newLimiter(time.Minute)
	if !l.Allow("lease-abc") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	l := newLimiter(time.Minute)
	l.Allow("lease-abc")
	if l.Allow("lease-abc") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_DifferentKeysIndependent(t *testing.T) {
	l := newLimiter(time.Minute)
	l.Allow("lease-a")
	if !l.Allow("lease-b") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestAllow_AllowedAfterCooldown(t *testing.T) {
	l := newLimiter(50 * time.Millisecond)
	l.Allow("lease-x")
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("lease-x") {
		t.Fatal("expected call to be allowed after cooldown elapsed")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	l := newLimiter(time.Minute)
	l.Allow("lease-r")
	l.Reset("lease-r")
	if !l.Allow("lease-r") {
		t.Fatal("expected allow after reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	l := newLimiter(time.Minute)
	l.Allow("k1")
	l.Allow("k2")
	if got := l.Len(); got != 2 {
		t.Fatalf("expected 2 keys, got %d", got)
	}
}

func TestPurge_RemovesExpiredKeys(t *testing.T) {
	l := newLimiter(50 * time.Millisecond)
	l.Allow("old")
	time.Sleep(60 * time.Millisecond)
	l.Allow("fresh")
	l.Purge()
	if got := l.Len(); got != 1 {
		t.Fatalf("expected 1 key after purge, got %d", got)
	}
}

func TestNew_ZeroCooldownDefaultsToMinute(t *testing.T) {
	l := ratelimit.New(0)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	// Two immediate calls — second must be blocked (default 1-minute cooldown).
	l.Allow("z")
	if l.Allow("z") {
		t.Fatal("expected second call to be blocked under default cooldown")
	}
}
