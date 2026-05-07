package cooldown

import (
	"testing"
	"time"
)

func newCooldown(period time.Duration) *Cooldown {
	c := New(period)
	return c
}

func TestAllow_NoRecord_ReturnsTrue(t *testing.T) {
	c := newCooldown(time.Minute)
	if !c.Allow("lease-1") {
		t.Fatal("expected Allow to return true for unseen lease")
	}
}

func TestAllow_WithinCooldown_ReturnsFalse(t *testing.T) {
	c := newCooldown(time.Minute)
	c.Record("lease-1")
	if c.Allow("lease-1") {
		t.Fatal("expected Allow to return false within cooldown period")
	}
}

func TestAllow_AfterCooldown_ReturnsTrue(t *testing.T) {
	c := newCooldown(50 * time.Millisecond)
	c.Record("lease-1")
	time.Sleep(60 * time.Millisecond)
	if !c.Allow("lease-1") {
		t.Fatal("expected Allow to return true after cooldown elapsed")
	}
}

func TestAllow_DifferentLeases_Independent(t *testing.T) {
	c := newCooldown(time.Minute)
	c.Record("lease-1")
	if !c.Allow("lease-2") {
		t.Fatal("expected lease-2 to be unaffected by lease-1 cooldown")
	}
}

func TestReset_ClearsCooldown(t *testing.T) {
	c := newCooldown(time.Minute)
	c.Record("lease-1")
	c.Reset("lease-1")
	if !c.Allow("lease-1") {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestPurge_RemovesExpiredRecords(t *testing.T) {
	c := newCooldown(50 * time.Millisecond)
	c.Record("lease-1")
	c.Record("lease-2")
	time.Sleep(60 * time.Millisecond)
	c.Purge()
	c.mu.Lock()
	n := len(c.records)
	c.mu.Unlock()
	if n != 0 {
		t.Fatalf("expected 0 records after Purge, got %d", n)
	}
}

func TestPurge_RetainsActiveCooldowns(t *testing.T) {
	c := newCooldown(time.Minute)
	c.Record("lease-active")
	c.Purge()
	if c.Allow("lease-active") {
		t.Fatal("expected active cooldown to survive Purge")
	}
}

func TestNew_DefaultPeriodUsedWhenZero(t *testing.T) {
	c := New(0)
	if c.period != DefaultCooldown {
		t.Fatalf("expected default period %v, got %v", DefaultCooldown, c.period)
	}
}
