package escalation_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/escalation"
	"github.com/your-org/vaultwatch/internal/lease"
)

func newEscalator() *escalation.Escalator {
	return escalation.New(escalation.Config{
		WarningAfter:   100 * time.Millisecond,
		CriticalAfter:  200 * time.Millisecond,
		EmergencyAfter: 300 * time.Millisecond,
	})
}

func sampleInfo(status lease.Status) lease.Info {
	return lease.Info{
		LeaseID: "secret/data/myapp/db#abc123",
		Status:  status,
	}
}

func TestEvaluate_HealthyReturnsNone(t *testing.T) {
	esc := newEscalator()
	tier := esc.Evaluate(sampleInfo(lease.StatusOK))
	if tier != escalation.TierNone {
		t.Fatalf("expected TierNone, got %d", tier)
	}
}

func TestEvaluate_ImmediatelyBelowWarning(t *testing.T) {
	esc := newEscalator()
	tier := esc.Evaluate(sampleInfo(lease.StatusWarning))
	if tier != escalation.TierNone {
		t.Fatalf("expected TierNone before warning threshold, got %d", tier)
	}
}

func TestEvaluate_ReachesWarningTier(t *testing.T) {
	esc := newEscalator()
	info := sampleInfo(lease.StatusWarning)
	esc.Evaluate(info) // seed entry
	time.Sleep(110 * time.Millisecond)
	tier := esc.Evaluate(info)
	if tier != escalation.TierWarning {
		t.Fatalf("expected TierWarning, got %d", tier)
	}
}

func TestEvaluate_ReachesCriticalTier(t *testing.T) {
	esc := newEscalator()
	info := sampleInfo(lease.StatusCritical)
	esc.Evaluate(info)
	time.Sleep(210 * time.Millisecond)
	tier := esc.Evaluate(info)
	if tier != escalation.TierCritical {
		t.Fatalf("expected TierCritical, got %d", tier)
	}
}

func TestEvaluate_ReachesEmergencyTier(t *testing.T) {
	esc := newEscalator()
	info := sampleInfo(lease.StatusCritical)
	esc.Evaluate(info)
	time.Sleep(310 * time.Millisecond)
	tier := esc.Evaluate(info)
	if tier != escalation.TierEmergency {
		t.Fatalf("expected TierEmergency, got %d", tier)
	}
}

func TestEvaluate_HealthyClearsState(t *testing.T) {
	esc := newEscalator()
	info := sampleInfo(lease.StatusCritical)
	esc.Evaluate(info)
	time.Sleep(210 * time.Millisecond)

	// lease recovers
	healthy := sampleInfo(lease.StatusOK)
	esc.Evaluate(healthy)

	// degrade again — timer should reset
	tier := esc.Evaluate(info)
	if tier != escalation.TierNone {
		t.Fatalf("expected TierNone after recovery, got %d", tier)
	}
}

func TestReset_ClearsTrackedLease(t *testing.T) {
	esc := newEscalator()
	info := sampleInfo(lease.StatusWarning)
	esc.Evaluate(info)
	time.Sleep(110 * time.Millisecond)
	esc.Reset(info.LeaseID)
	tier := esc.Evaluate(info)
	if tier != escalation.TierNone {
		t.Fatalf("expected TierNone after reset, got %d", tier)
	}
}
