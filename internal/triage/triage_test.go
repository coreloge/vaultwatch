package triage_test

import (
	"testing"
	"time"

	"github.com/youorg/vaultwatch/internal/lease"
	"github.com/youorg/vaultwatch/internal/triage"
)

func newInfo(status lease.Status) lease.Info {
	return lease.Info{
		LeaseID:   "secret/data/test",
		Status:    status,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func TestNew_ReturnsEmptyQueue(t *testing.T) {
	q := triage.New()
	if q.Len() != 0 {
		t.Fatalf("expected empty queue, got len %d", q.Len())
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	q := triage.New()
	q.Add(newInfo(lease.StatusWarning))
	q.Add(newInfo(lease.StatusCritical))
	if q.Len() != 2 {
		t.Fatalf("expected len 2, got %d", q.Len())
	}
}

func TestDrain_ClearsQueue(t *testing.T) {
	q := triage.New()
	q.Add(newInfo(lease.StatusWarning))
	q.Drain()
	if q.Len() != 0 {
		t.Fatalf("expected empty queue after drain, got %d", q.Len())
	}
}

func TestDrain_OrdersByPriorityDescending(t *testing.T) {
	q := triage.New()
	q.Add(newInfo(lease.StatusOK))
	q.Add(newInfo(lease.StatusWarning))
	q.Add(newInfo(lease.StatusCritical))

	entries := q.Drain()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Priority != triage.PriorityHigh {
		t.Errorf("expected first entry to be PriorityHigh, got %v", entries[0].Priority)
	}
	if entries[1].Priority != triage.PriorityMedium {
		t.Errorf("expected second entry to be PriorityMedium, got %v", entries[1].Priority)
	}
	if entries[2].Priority != triage.PriorityLow {
		t.Errorf("expected third entry to be PriorityLow, got %v", entries[2].Priority)
	}
}

func TestDrain_SamePriorityOrderedByAge(t *testing.T) {
	q := triage.New()
	old := lease.Info{
		LeaseID:   "secret/old",
		Status:    lease.StatusWarning,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	new_ := lease.Info{
		LeaseID:   "secret/new",
		Status:    lease.StatusWarning,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	q.Add(old)
	time.Sleep(2 * time.Millisecond)
	q.Add(new_)

	entries := q.Drain()
	if entries[0].Info.LeaseID != "secret/old" {
		t.Errorf("expected older entry first, got %s", entries[0].Info.LeaseID)
	}
}

func TestDrain_ExpiredIsHighPriority(t *testing.T) {
	q := triage.New()
	q.Add(newInfo(lease.StatusExpired))
	entries := q.Drain()
	if entries[0].Priority != triage.PriorityHigh {
		t.Errorf("expected expired lease to be PriorityHigh, got %v", entries[0].Priority)
	}
}
