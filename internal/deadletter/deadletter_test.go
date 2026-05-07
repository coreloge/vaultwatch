package deadletter_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/deadletter"
	"github.com/your-org/vaultwatch/internal/lease"
)

func sampleInfo(id string) lease.Info {
	return lease.Info{LeaseID: id, Path: "secret/data/" + id}
}

func newStore() *deadletter.Store {
	return deadletter.New(8, time.Minute)
}

func TestAdd_IncreasesLen(t *testing.T) {
	s := newStore()
	s.Add(sampleInfo("lease-1"), "timeout", 3)
	if got := s.Len(); got != 1 {
		t.Fatalf("expected len 1, got %d", got)
	}
}

func TestDrain_ReturnsEntries(t *testing.T) {
	s := newStore()
	s.Add(sampleInfo("lease-1"), "timeout", 3)
	s.Add(sampleInfo("lease-2"), "connect refused", 5)

	entries := s.Drain()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].LeaseInfo.LeaseID != "lease-1" {
		t.Errorf("unexpected first entry: %s", entries[0].LeaseInfo.LeaseID)
	}
}

func TestDrain_ClearsStore(t *testing.T) {
	s := newStore()
	s.Add(sampleInfo("lease-1"), "timeout", 1)
	s.Drain()
	if got := s.Len(); got != 0 {
		t.Fatalf("expected empty store after drain, got %d", got)
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	s := deadletter.New(3, time.Minute)
	for i := 0; i < 4; i++ {
		s.Add(sampleInfo("lease-"+string(rune('a'+i))), "err", 1)
	}
	if got := s.Len(); got != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", got)
	}
	entries := s.Drain()
	// The first entry (lease-a) should have been evicted.
	if entries[0].LeaseInfo.LeaseID == "lease-a" {
		t.Error("expected oldest entry to be evicted")
	}
}

func TestAdd_ExpiredEntriesPurged(t *testing.T) {
	s := deadletter.New(8, time.Millisecond)
	s.Add(sampleInfo("lease-1"), "timeout", 2)
	time.Sleep(5 * time.Millisecond)
	// Adding a new entry triggers purge of expired ones.
	s.Add(sampleInfo("lease-2"), "timeout", 1)
	if got := s.Len(); got != 1 {
		t.Fatalf("expected 1 live entry, got %d", got)
	}
}

func TestEntry_FieldsPreserved(t *testing.T) {
	s := newStore()
	info := sampleInfo("lease-xyz")
	s.Add(info, "connection reset", 7)
	entries := s.Drain()
	e := entries[0]
	if e.Reason != "connection reset" {
		t.Errorf("reason mismatch: %s", e.Reason)
	}
	if e.Attempts != 7 {
		t.Errorf("attempts mismatch: %d", e.Attempts)
	}
	if e.FailedAt.IsZero() {
		t.Error("FailedAt should not be zero")
	}
	if e.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should not be zero")
	}
}

func TestNew_DefaultsApplied(t *testing.T) {
	s := deadletter.New(0, 0)
	if s == nil {
		t.Fatal("expected non-nil store with zero defaults")
	}
}
