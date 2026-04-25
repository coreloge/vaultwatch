package filter_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/filter"
	"github.com/yourusername/vaultwatch/internal/lease"
)

func newInfo(id string, status lease.Status) lease.Info {
	return lease.Info{
		LeaseID:   id,
		Status:    status,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func TestAllow_NoRules_AllowsAll(t *testing.T) {
	f := filter.New(nil, nil)
	info := newInfo("secret/data/myapp/token", lease.StatusWarning)
	if !f.Allow(info) {
		t.Error("expected lease to be allowed with no rules")
	}
}

func TestAllow_ExcludeByPrefix(t *testing.T) {
	excludes := []filter.Rule{{PathPrefix: "secret/data/myapp"}}
	f := filter.New(nil, excludes)

	if f.Allow(newInfo("secret/data/myapp/token", lease.StatusOK)) {
		t.Error("expected lease to be excluded by prefix")
	}
	if !f.Allow(newInfo("secret/data/otherapp/token", lease.StatusOK)) {
		t.Error("expected other lease to be allowed")
	}
}

func TestAllow_IncludeByStatus(t *testing.T) {
	includes := []filter.Rule{{Statuses: []lease.Status{lease.StatusCritical, lease.StatusExpired}}}
	f := filter.New(includes, nil)

	if !f.Allow(newInfo("secret/a", lease.StatusCritical)) {
		t.Error("expected critical lease to be included")
	}
	if f.Allow(newInfo("secret/b", lease.StatusOK)) {
		t.Error("expected OK lease to be excluded when not in include list")
	}
}

func TestAllow_ExcludeTakesPrecedence(t *testing.T) {
	includes := []filter.Rule{{PathPrefix: "secret/"}}
	excludes := []filter.Rule{{PathPrefix: "secret/data/sensitive"}}
	f := filter.New(includes, excludes)

	if f.Allow(newInfo("secret/data/sensitive/key", lease.StatusWarning)) {
		t.Error("expected excluded lease to be blocked even if it matches include")
	}
	if !f.Allow(newInfo("secret/data/safe/key", lease.StatusWarning)) {
		t.Error("expected non-excluded lease under include prefix to be allowed")
	}
}

func TestAllow_IncludeByPrefixAndStatus(t *testing.T) {
	includes := []filter.Rule{{
		PathPrefix: "aws/",
		Statuses:   []lease.Status{lease.StatusWarning},
	}}
	f := filter.New(includes, nil)

	if !f.Allow(newInfo("aws/creds/role", lease.StatusWarning)) {
		t.Error("expected matching prefix+status to be allowed")
	}
	if f.Allow(newInfo("aws/creds/role", lease.StatusOK)) {
		t.Error("expected wrong status to be blocked")
	}
	if f.Allow(newInfo("gcp/creds/role", lease.StatusWarning)) {
		t.Error("expected wrong prefix to be blocked")
	}
}
