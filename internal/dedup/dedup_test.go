package dedup_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/dedup"
	"github.com/your-org/vaultwatch/internal/lease"
)

func newDeduplicator(window time.Duration) *dedup.Deduplicator {
	return dedup.New(window)
}

func TestIsDuplicate_NotRecorded(t *testing.T) {
	dd := newDeduplicator(time.Minute)
	if dd.IsDuplicate("lease-1", lease.StatusWarning) {
		t.Fatal("expected false for unseen lease")
	}
}

func TestIsDuplicate_SameStatusWithinWindow(t *testing.T) {
	dd := newDeduplicator(time.Minute)
	dd.Record("lease-1", lease.StatusWarning)
	if !dd.IsDuplicate("lease-1", lease.StatusWarning) {
		t.Fatal("expected duplicate within window")
	}
}

func TestIsDuplicate_StatusChanged(t *testing.T) {
	dd := newDeduplicator(time.Minute)
	dd.Record("lease-1", lease.StatusWarning)
	if dd.IsDuplicate("lease-1", lease.StatusCritical) {
		t.Fatal("status changed — should not be duplicate")
	}
}

func TestIsDuplicate_WindowExpired(t *testing.T) {
	dd := newDeduplicator(10 * time.Millisecond)
	dd.Record("lease-1", lease.StatusWarning)
	time.Sleep(20 * time.Millisecond)
	if dd.IsDuplicate("lease-1", lease.StatusWarning) {
		t.Fatal("window expired — should not be duplicate")
	}
}

func TestEvict_ClearsEntry(t *testing.T) {
	dd := newDeduplicator(time.Minute)
	dd.Record("lease-1", lease.StatusCritical)
	dd.Evict("lease-1")
	if dd.IsDuplicate("lease-1", lease.StatusCritical) {
		t.Fatal("entry should have been evicted")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	dd := newDeduplicator(10 * time.Millisecond)
	dd.Record("lease-1", lease.StatusWarning)
	dd.Record("lease-2", lease.StatusCritical)
	time.Sleep(20 * time.Millisecond)
	dd.Purge()
	if dd.IsDuplicate("lease-1", lease.StatusWarning) {
		t.Fatal("lease-1 should have been purged")
	}
	if dd.IsDuplicate("lease-2", lease.StatusCritical) {
		t.Fatal("lease-2 should have been purged")
	}
}

func TestPurge_KeepsActiveEntries(t *testing.T) {
	dd := newDeduplicator(time.Minute)
	dd.Record("lease-1", lease.StatusWarning)
	dd.Purge()
	if !dd.IsDuplicate("lease-1", lease.StatusWarning) {
		t.Fatal("active entry should survive purge")
	}
}
