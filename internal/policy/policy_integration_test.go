package policy_test

import (
	"sync"
	"testing"

	"github.com/your-org/vaultwatch/internal/lease"
	"github.com/your-org/vaultwatch/internal/policy"
)

// TestConcurrentEvaluate_NoPanic verifies that a shared Policy is safe to call
// from multiple goroutines simultaneously (it is read-only after construction).
func TestConcurrentEvaluate_NoPanic(t *testing.T) {
	p := policy.New([]policy.Rule{
		{PathPrefix: "secret/", Action: policy.ActionSuppress},
		{Statuses: []lease.Status{lease.StatusCritical}, Action: policy.ActionEscalate},
	})

	infos := []lease.Info{
		newInfo("secret/db", lease.StatusWarning),
		newInfo("auth/token", lease.StatusCritical),
		newInfo("pki/cert", lease.StatusOK),
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			info := infos[i%len(infos)]
			p.Evaluate(info) // must not panic
		}(i)
	}
	wg.Wait()
}

// TestEvaluate_LargeRuleSet ensures correct first-match behaviour across many rules.
func TestEvaluate_LargeRuleSet(t *testing.T) {
	rules := make([]policy.Rule, 20)
	for i := range rules {
		rules[i] = policy.Rule{
			PathPrefix: "nomatch/",
			Action:     policy.ActionSuppress,
		}
	}
	// Add a matching rule at the end.
	rules = append(rules, policy.Rule{
		PathPrefix: "secret/",
		Action:     policy.ActionEscalate,
	})

	p := policy.New(rules)
	action, _ := p.Evaluate(newInfo("secret/foo", lease.StatusWarning))
	if action != policy.ActionEscalate {
		t.Fatalf("expected escalate from last rule, got %s", action)
	}
}
