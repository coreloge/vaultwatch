package lease_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
)

func sampleLease(id string) lease.Info {
	return lease.Info{
		LeaseID:   id,
		Path:      "secret/data/app",
		Renewable: true,
		TTL:       1 * time.Hour,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
}

func TestStore_SetAndGet(t *testing.T) {
	s := lease.NewStore()
	l := sampleLease("lease/abc")
	s.Set(l)

	got, err := s.Get("lease/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.LeaseID != l.LeaseID {
		t.Errorf("expected %q, got %q", l.LeaseID, got.LeaseID)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	s := lease.NewStore()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing lease")
	}
}

func TestStore_Delete(t *testing.T) {
	s := lease.NewStore()
	s.Set(sampleLease("lease/del"))
	s.Delete("lease/del")
	_, err := s.Get("lease/del")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestStore_All(t *testing.T) {
	s := lease.NewStore()
	s.Set(sampleLease("lease/1"))
	s.Set(sampleLease("lease/2"))
	s.Set(sampleLease("lease/3"))

	all := s.All()
	if len(all) != 3 {
		t.Errorf("expected 3 leases, got %d", len(all))
	}
}

func TestStore_Count(t *testing.T) {
	s := lease.NewStore()
	if s.Count() != 0 {
		t.Error("expected empty store")
	}
	s.Set(sampleLease("lease/x"))
	if s.Count() != 1 {
		t.Errorf("expected count 1, got %d", s.Count())
	}
}

func TestStore_UpdateExisting(t *testing.T) {
	s := lease.NewStore()
	l := sampleLease("lease/upd")
	s.Set(l)

	updated := l
	updated.TTL = 2 * time.Hour
	updated.ExpiresAt = time.Now().Add(2 * time.Hour)
	s.Set(updated)

	got, _ := s.Get("lease/upd")
	if got.TTL != 2*time.Hour {
		t.Errorf("expected updated TTL, got %v", got.TTL)
	}
}
