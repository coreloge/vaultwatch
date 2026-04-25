package suppress_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/suppress"
)

func newSuppressor(d time.Duration) *suppress.Suppressor {
	return suppress.New(d)
}

func TestIsSuppressed_NotSuppressed(t *testing.T) {
	s := newSuppressor(time.Minute)
	if s.IsSuppressed("lease-abc") {
		t.Fatal("expected lease to not be suppressed")
	}
}

func TestSuppress_ThenIsSuppressed(t *testing.T) {
	s := newSuppressor(time.Minute)
	s.Suppress("lease-abc")
	if !s.IsSuppressed("lease-abc") {
		t.Fatal("expected lease to be suppressed after Suppress()")
	}
}

func TestIsSuppressed_Expired(t *testing.T) {
	s := newSuppressor(10 * time.Millisecond)
	s.Suppress("lease-xyz")
	time.Sleep(30 * time.Millisecond)
	if s.IsSuppressed("lease-xyz") {
		t.Fatal("expected suppression window to have expired")
	}
}

func TestRelease_ClearsSuppression(t *testing.T) {
	s := newSuppressor(time.Minute)
	s.Suppress("lease-123")
	s.Release("lease-123")
	if s.IsSuppressed("lease-123") {
		t.Fatal("expected lease to be released")
	}
}

func TestActive_CountsWindows(t *testing.T) {
	s := newSuppressor(time.Minute)
	s.Suppress("a")
	s.Suppress("b")
	s.Suppress("c")
	if got := s.Active(); got != 3 {
		t.Fatalf("expected 3 active windows, got %d", got)
	}
}

func TestActive_ExcludesExpired(t *testing.T) {
	s := newSuppressor(10 * time.Millisecond)
	s.Suppress("expiring")
	s.Suppress("expiring2")
	time.Sleep(30 * time.Millisecond)

	// Add a long-lived one
	s2 := suppress.New(time.Minute)
	s2.Suppress("permanent")

	if got := s.Active(); got != 0 {
		t.Fatalf("expected 0 active windows after expiry, got %d", got)
	}
	if got := s2.Active(); got != 1 {
		t.Fatalf("expected 1 active window, got %d", got)
	}
}

func TestSuppress_IndependentLeases(t *testing.T) {
	s := newSuppressor(time.Minute)
	s.Suppress("lease-1")
	if s.IsSuppressed("lease-2") {
		t.Fatal("lease-2 should not be suppressed")
	}
}
