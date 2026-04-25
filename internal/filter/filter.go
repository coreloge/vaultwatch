// Package filter provides lease filtering capabilities for vaultwatch,
// allowing operators to include or exclude leases by path prefix or status.
package filter

import (
	"strings"

	"github.com/yourusername/vaultwatch/internal/lease"
)

// Rule defines a single include/exclude filter rule.
type Rule struct {
	PathPrefix string
	Statuses   []lease.Status
}

// Filter evaluates leases against a set of include and exclude rules.
type Filter struct {
	includes []Rule
	excludes []Rule
}

// New creates a Filter with the given include and exclude rules.
// If no include rules are provided, all leases are included by default.
func New(includes, excludes []Rule) *Filter {
	return &Filter{
		includes: includes,
		excludes: excludes,
	}
}

// Allow returns true if the lease info passes the filter.
func (f *Filter) Allow(info lease.Info) bool {
	if f.matchesAny(info, f.excludes) {
		return false
	}
	if len(f.includes) == 0 {
		return true
	}
	return f.matchesAny(info, f.includes)
}

func (f *Filter) matchesAny(info lease.Info, rules []Rule) bool {
	for _, r := range rules {
		if matchesRule(info, r) {
			return true
		}
	}
	return false
}

func matchesRule(info lease.Info, r Rule) bool {
	if r.PathPrefix != "" && !strings.HasPrefix(info.LeaseID, r.PathPrefix) {
		return false
	}
	if len(r.Statuses) > 0 {
		for _, s := range r.Statuses {
			if info.Status == s {
				return true
			}
		}
		return false
	}
	return true
}
