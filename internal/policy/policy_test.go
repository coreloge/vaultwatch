package policy_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
	"github.com/your-org/vaultwatch/internal/policy"
)

func newInfo(id string, status lease.Status) lease.Info {
	return lease.Info{
		LeaseID: id,
		Status:  status,
	}
}

func TestEvaluate_NoRules_AllowsAll(t *testing.T) {
	p := policy.New(nil)
	action, dur := p.Evaluate(newInfo("secret/foo", lease.StatusWarning))
	if action != policy.ActionAllow {
		t.Fatalf("expected allow, got %s", action)
	}
	if dur != 0 {
		t.Fatalf("expected zero duration, got %s", dur)
	}
}

func TestEvaluate_PrefixMatch_Suppresses(t *testing.T) {
	p := policy.New([]policy.Rule{
		{PathPrefix: "secret/internal/", Action: policy.ActionSuppress, SuppressDuration: 10 * time.Minute},
	})
	action, dur := p.Evaluate(newInfo("secret/internal/db", lease.StatusCritical))
	if action != policy.ActionSuppress {
		t.Fatalf("expected suppress, got %s", action)
	}
	if dur != 10*time.Minute {
		t.Fatalf("expected 10m, got %s", dur)
	}
}

func TestEvaluate_PrefixNoMatch_Falls Through(t *testing.T) {
	p := policy.New([]policy.Rule{
		{PathPrefix: "secret/internal/", Action: policy.ActionSuppress},
	})
	action, _ := p.Evaluate(newInfo("secret/public/cert", lease.StatusWarning))
	if action != policy.ActionAllow {
		t.Fatalf("expected allow, got %s", action)
	}
}

func TestEvaluate_StatusMatch_Escalates(t *testing.T) {
	p := policy.New([]policy.Rule{
		{Statuses: []lease.Status{lease.StatusCritical}, Action: policy.ActionEscalate},
	})
	action, _ := p.Evaluate(newInfo("secret/foo", lease.StatusCritical))
	if action != policy.ActionEscalate {
		t.Fatalf("expected escalate, got %s", action)
	}
}

func TestEvaluate_StatusNoMatch_Allows(t *testing.T) {
	p := policy.New([]policy.Rule{
		{Statuses: []lease.Status{lease.StatusCritical}, Action: policy.ActionEscalate},
	})
	action, _ := p.Evaluate(newInfo("secret/foo", lease.StatusWarning))
	if action != policy.ActionAllow {
		t.Fatalf("expected allow, got %s", action)
	}
}

func TestEvaluate_FirstRuleWins(t *testing.T) {
	p := policy.New([]policy.Rule{
		{PathPrefix: "secret/", Action: policy.ActionSuppress},
		{PathPrefix: "secret/critical/", Action: policy.ActionEscalate},
	})
	// Both rules match, but first one should win.
	action, _ := p.Evaluate(newInfo("secret/critical/token", lease.StatusCritical))
	if action != policy.ActionSuppress {
		t.Fatalf("expected suppress (first rule), got %s", action)
	}
}

func TestEvaluate_CombinedPrefixAndStatus(t *testing.T) {
	p := policy.New([]policy.Rule{
		{
			PathPrefix: "auth/",
			Statuses:   []lease.Status{lease.StatusWarning},
			Action:     policy.ActionSuppress,
		},
	})
	// Prefix matches but status does not.
	action, _ := p.Evaluate(newInfo("auth/token/abc", lease.StatusCritical))
	if action != policy.ActionAllow {
		t.Fatalf("expected allow (status mismatch), got %s", action)
	}
	// Both match.
	action, _ = p.Evaluate(newInfo("auth/token/abc", lease.StatusWarning))
	if action != policy.ActionSuppress {
		t.Fatalf("expected suppress, got %s", action)
	}
}
