// Package routing provides lease-aware alert routing, directing alerts to
// one or more webhook destinations based on configurable routing rules.
// Rules are matched against lease metadata such as path prefix, status,
// and mount point, enabling fine-grained delivery control.
package routing

import (
	"strings"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Rule defines a single routing rule that maps matching leases to a set
// of webhook destinations.
type Rule struct {
	// PathPrefix restricts the rule to leases whose path begins with this
	// value. An empty string matches all paths.
	PathPrefix string

	// Statuses restricts the rule to leases with one of the listed statuses.
	// An empty slice matches all statuses.
	Statuses []lease.Status

	// Destinations is the list of webhook URLs that should receive alerts
	// when this rule matches.
	Destinations []string
}

// Router maps lease info to a set of webhook destinations using an ordered
// list of rules. The first matching rule wins; if no rule matches, the
// default destinations are used.
type Router struct {
	rules        []Rule
	defaultDests []string
}

// New creates a Router with the provided rules and default destinations.
// Rules are evaluated in the order supplied; the first match is used.
// defaultDests is returned when no rule matches.
func New(rules []Rule, defaultDests []string) *Router {
	return &Router{
		rules:        rules,
		defaultDests: defaultDests,
	}
}

// Route returns the webhook destinations for the given lease info.
// It evaluates each rule in order and returns the destinations of the
// first matching rule. If no rule matches, the router's default
// destinations are returned. An empty slice is returned only when both
// no rule matches and no defaults are configured.
func (r *Router) Route(info lease.Info) []string {
	for _, rule := range r.rules {
		if r.matches(rule, info) {
			return rule.Destinations
		}
	}
	return r.defaultDests
}

// Rules returns a copy of the configured routing rules.
func (r *Router) Rules() []Rule {
	out := make([]Rule, len(r.rules))
	copy(out, r.rules)
	return out
}

// DefaultDestinations returns the destinations used when no rule matches.
func (r *Router) DefaultDestinations() []string {
	out := make([]string, len(r.defaultDests))
	copy(out, r.defaultDests)
	return out
}

// matches reports whether rule applies to the given lease info.
func (r *Router) matches(rule Rule, info lease.Info) bool {
	if rule.PathPrefix != "" && !strings.HasPrefix(info.LeaseID, rule.PathPrefix) {
		return false
	}
	if len(rule.Statuses) > 0 && !statusIn(info.Status, rule.Statuses) {
		return false
	}
	return true
}

// statusIn reports whether s is present in the statuses slice.
func statusIn(s lease.Status, statuses []lease.Status) bool {
	for _, allowed := range statuses {
		if s == allowed {
			return true
		}
	}
	return false
}
