// Package policy defines alert suppression and routing policies
// based on lease metadata, status, and configurable rules.
package policy

import (
	"strings"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Action describes what should happen when a policy matches.
type Action string

const (
	ActionAllow    Action = "allow"
	ActionSuppress Action = "suppress"
	ActionEscalate Action = "escalate"
)

// Rule defines a single matching rule and its resulting action.
type Rule struct {
	// PathPrefix restricts the rule to leases whose ID starts with this prefix.
	PathPrefix string
	// Statuses restricts the rule to specific lease statuses. Empty means all.
	Statuses []lease.Status
	// Action is the outcome when this rule matches.
	Action Action
	// SuppressDuration is only used when Action == ActionSuppress.
	SuppressDuration time.Duration
}

// Policy evaluates a set of ordered rules against a lease.
type Policy struct {
	rules []Rule
}

// New creates a Policy with the given rules. Rules are evaluated in order;
// the first match wins. If no rule matches, ActionAllow is returned.
func New(rules []Rule) *Policy {
	return &Policy{rules: rules}
}

// Evaluate returns the Action for the given lease info.
func (p *Policy) Evaluate(info lease.Info) (Action, time.Duration) {
	for _, r := range p.rules {
		if !p.matches(r, info) {
			continue
		}
		return r.Action, r.SuppressDuration
	}
	return ActionAllow, 0
}

func (p *Policy) matches(r Rule, info lease.Info) bool {
	if r.PathPrefix != "" && !strings.HasPrefix(info.LeaseID, r.PathPrefix) {
		return false
	}
	if len(r.Statuses) > 0 && !statusIn(info.Status, r.Statuses) {
		return false
	}
	return true
}

func statusIn(s lease.Status, list []lease.Status) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
