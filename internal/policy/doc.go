// Package policy provides a rule-based evaluation engine for VaultWatch lease
// alerts. Each [Rule] matches leases by path prefix and/or status, and assigns
// an [Action] of allow, suppress, or escalate.
//
// Rules are evaluated in declaration order; the first matching rule wins.
// When no rule matches the default action is [ActionAllow].
//
// Example:
//
//	p := policy.New([]policy.Rule{
//		{PathPrefix: "secret/internal/", Action: policy.ActionSuppress},
//		{Statuses: []lease.Status{lease.StatusCritical}, Action: policy.ActionEscalate},
//	})
//	action, _ := p.Evaluate(info)
package policy
