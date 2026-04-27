package envelope_test

import (
	"testing"
	"time"

	"github.com/youorg/vaultwatch/internal/envelope"
	"github.com/youorg/vaultwatch/internal/lease"
)

func sampleInfo() lease.Info {
	return lease.Info{
		LeaseID:   "secret/data/db#abc123",
		Renewable: true,
		TTL:       lease.NewTTLFromSeconds(300),
	}
}

func TestNew_AssignsID(t *testing.T) {
	e := envelope.New(sampleInfo(), "monitor")
	if e.ID == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestNew_SetsOrigin(t *testing.T) {
	e := envelope.New(sampleInfo(), "watcher")
	if e.Origin != "watcher" {
		t.Fatalf("expected origin 'watcher', got %q", e.Origin)
	}
}

func TestNew_AttemptsZero(t *testing.T) {
	e := envelope.New(sampleInfo(), "monitor")
	if e.Attempts != 0 {
		t.Fatalf("expected 0 attempts, got %d", e.Attempts)
	}
}

func TestIncrement_IncreasesAttempts(t *testing.T) {
	e := envelope.New(sampleInfo(), "monitor")
	e.Increment()
	e.Increment()
	if e.Attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", e.Attempts)
	}
}

func TestAge_NonNegative(t *testing.T) {
	e := envelope.New(sampleInfo(), "monitor")
	time.Sleep(2 * time.Millisecond)
	if e.Age() <= 0 {
		t.Fatal("expected positive age")
	}
}

func TestString_ContainsLeaseID(t *testing.T) {
	info := sampleInfo()
	e := envelope.New(info, "monitor")
	s := e.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
	if !containsSubstring(s, info.LeaseID) {
		t.Fatalf("expected string to contain lease ID %q, got %q", info.LeaseID, s)
	}
}

func TestNew_UniqueIDs(t *testing.T) {
	a := envelope.New(sampleInfo(), "monitor")
	b := envelope.New(sampleInfo(), "monitor")
	if a.ID == b.ID {
		t.Fatal("expected unique IDs for distinct envelopes")
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
