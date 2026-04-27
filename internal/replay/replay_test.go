package replay_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
	"github.com/your-org/vaultwatch/internal/replay"
)

func sampleInfo(id string) lease.Info {
	return lease.Info{
		LeaseID:   id,
		MountPath: "secret/",
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	s := replay.New(time.Hour)
	s.Add(sampleInfo("lease-1"))
	s.Add(sampleInfo("lease-2"))
	if got := s.Len(); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestDrain_ReturnsAllEntries(t *testing.T) {
	s := replay.New(time.Hour)
	s.Add(sampleInfo("lease-a"))
	s.Add(sampleInfo("lease-b"))

	entries := s.Drain()
	if len(entries) != 2 {
		t.Fatalf("expected 2 drained entries, got %d", len(entries))
	}
	if s.Len() != 0 {
		t.Errorf("expected store to be empty after drain, got %d", s.Len())
	}
}

func TestDrain_DropsStalEntries(t *testing.T) {
	s := replay.New(1 * time.Millisecond)
	s.Add(sampleInfo("stale-lease"))

	time.Sleep(5 * time.Millisecond)

	entries := s.Drain()
	if len(entries) != 0 {
		t.Errorf("expected stale entry to be discarded, got %d entries", len(entries))
	}
}

func TestDrain_EmptyStore(t *testing.T) {
	s := replay.New(time.Hour)
	entries := s.Drain()
	if entries != nil && len(entries) != 0 {
		t.Errorf("expected nil or empty slice, got %v", entries)
	}
}

func TestPurge_ClearsAllEntries(t *testing.T) {
	s := replay.New(time.Hour)
	s.Add(sampleInfo("lease-x"))
	s.Add(sampleInfo("lease-y"))
	s.Purge()
	if s.Len() != 0 {
		t.Errorf("expected 0 entries after purge, got %d", s.Len())
	}
}

func TestNew_DefaultMaxAge(t *testing.T) {
	// zero maxAge should default to 1 hour — store should accept entries
	s := replay.New(0)
	s.Add(sampleInfo("lease-default"))
	if s.Len() != 1 {
		t.Errorf("expected 1 entry with default maxAge, got %d", s.Len())
	}
}
