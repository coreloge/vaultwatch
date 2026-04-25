package snapshot_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
	"github.com/yourusername/vaultwatch/internal/snapshot"
)

func sampleInfo(id string, status lease.Status) lease.Info {
	return lease.Info{
		LeaseID:   id,
		Status:    status,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func TestNewStore_ReturnsNonNil(t *testing.T) {
	s := snapshot.NewStore()
	if s == nil {
		t.Fatal("expected non-nil Store")
	}
}

func TestCompare_NoPreviousSnapshot(t *testing.T) {
	s := snapshot.NewStore()
	info := sampleInfo("lease/abc", lease.StatusWarning)

	changed := s.Compare(info)
	if !changed {
		t.Error("expected changed=true when no prior snapshot exists")
	}
}

func TestCompare_SameStatus(t *testing.T) {
	s := snapshot.NewStore()
	info := sampleInfo("lease/abc", lease.StatusWarning)

	s.Compare(info) // record initial
	changed := s.Compare(info)
	if changed {
		t.Error("expected changed=false when status unchanged")
	}
}

func TestCompare_StatusChanged(t *testing.T) {
	s := snapshot.NewStore()

	first := sampleInfo("lease/abc", lease.StatusOK)
	s.Compare(first)

	second := sampleInfo("lease/abc", lease.StatusCritical)
	changed := s.Compare(second)
	if !changed {
		t.Error("expected changed=true when status transitions")
	}
}

func TestDelete_RemovesSnapshot(t *testing.T) {
	s := snapshot.NewStore()
	info := sampleInfo("lease/xyz", lease.StatusOK)

	s.Compare(info)
	s.Delete("lease/xyz")

	// After deletion, next compare should treat it as new
	changed := s.Compare(info)
	if !changed {
		t.Error("expected changed=true after deletion and re-compare")
	}
}

func TestAll_ReturnsAllSnapshots(t *testing.T) {
	s := snapshot.NewStore()

	s.Compare(sampleInfo("lease/a", lease.StatusOK))
	s.Compare(sampleInfo("lease/b", lease.StatusWarning))
	s.Compare(sampleInfo("lease/c", lease.StatusCritical))

	all := s.All()
	if len(all) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(all))
	}
}
