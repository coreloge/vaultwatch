package digest_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/digest"
	"github.com/your-org/vaultwatch/internal/lease"
)

func sampleInfo(leaseID, status string, expiresAt time.Time) lease.Info {
	return lease.Info{
		LeaseID:   leaseID,
		Status:    lease.Status(status),
		ExpiresAt: expiresAt,
	}
}

func TestCompute_ReturnsTruncatedLength(t *testing.T) {
	d := digest.New(16)
	info := sampleInfo("lease/abc", "warning", time.Now().Add(5*time.Minute))
	got := d.Compute(info)
	if len(got) != 16 {
		t.Fatalf("expected length 16, got %d", len(got))
	}
}

func TestCompute_FullLengthWhenZero(t *testing.T) {
	d := digest.New(0)
	info := sampleInfo("lease/abc", "warning", time.Now().Add(5*time.Minute))
	got := d.Compute(info)
	if len(got) != 64 {
		t.Fatalf("expected full SHA-256 hex length 64, got %d", len(got))
	}
}

func TestCompute_Deterministic(t *testing.T) {
	d := digest.New(32)
	expiry := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	info := sampleInfo("lease/xyz", "critical", expiry)
	a := d.Compute(info)
	b := d.Compute(info)
	if a != b {
		t.Fatalf("expected identical digests, got %q and %q", a, b)
	}
}

func TestCompute_DifferentLeaseIDs_DifferentDigests(t *testing.T) {
	d := digest.New(0)
	expiry := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	a := d.Compute(sampleInfo("lease/aaa", "critical", expiry))
	b := d.Compute(sampleInfo("lease/bbb", "critical", expiry))
	if a == b {
		t.Fatal("expected different digests for different lease IDs")
	}
}

func TestCompute_SameMinuteBoundary_SameDigest(t *testing.T) {
	d := digest.New(0)
	base := time.Date(2025, 6, 1, 12, 30, 0, 0, time.UTC)
	infoA := sampleInfo("lease/t", "warning", base.Add(10*time.Second))
	infoB := sampleInfo("lease/t", "warning", base.Add(45*time.Second))
	if d.Compute(infoA) != d.Compute(infoB) {
		t.Fatal("expected same digest within the same minute boundary")
	}
}

func TestCompute_AcrossMinuteBoundary_DifferentDigest(t *testing.T) {
	d := digest.New(0)
	base := time.Date(2025, 6, 1, 12, 30, 0, 0, time.UTC)
	infoA := sampleInfo("lease/t", "warning", base)
	infoB := sampleInfo("lease/t", "warning", base.Add(time.Minute))
	if d.Compute(infoA) == d.Compute(infoB) {
		t.Fatal("expected different digest across minute boundary")
	}
}

func TestEqual_IdenticalStrings(t *testing.T) {
	if !digest.Equal("abc123", "abc123") {
		t.Fatal("expected Equal to return true for identical strings")
	}
}

func TestEqual_DifferentStrings(t *testing.T) {
	if digest.Equal("abc123", "def456") {
		t.Fatal("expected Equal to return false for different strings")
	}
}
