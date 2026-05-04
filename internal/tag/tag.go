// Package tag provides lease tagging and label enrichment for alert payloads.
package tag

import (
	"strings"
	"sync"
)

// Tagger attaches static and dynamic tags to lease identifiers.
type Tagger struct {
	mu      sync.RWMutex
	static  map[string]string
	prefix  map[string]map[string]string // prefix -> tags
}

// New returns a Tagger with the given static tags applied to all leases.
func New(static map[string]string) *Tagger {
	copied := make(map[string]string, len(static))
	for k, v := range static {
		copied[k] = v
	}
	return &Tagger{
		static: copied,
		prefix: make(map[string]map[string]string),
	}
}

// AddPrefix registers tags that apply to all lease IDs with the given prefix.
func (t *Tagger) AddPrefix(prefix string, tags map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	copied := make(map[string]string, len(tags))
	for k, v := range tags {
		copied[k] = v
	}
	t.prefix[prefix] = copied
}

// Tag returns the merged tag set for the given lease ID.
// Static tags are applied first; prefix tags override statics; later-registered
// prefixes take precedence over earlier ones when both match.
func (t *Tagger) Tag(leaseID string) map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make(map[string]string, len(t.static))
	for k, v := range t.static {
		out[k] = v
	}
	for pfx, tags := range t.prefix {
		if strings.HasPrefix(leaseID, pfx) {
			for k, v := range tags {
				out[k] = v
			}
		}
	}
	return out
}

// Keys returns all tag keys currently registered (static + prefix).
func (t *Tagger) Keys() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	seen := make(map[string]struct{})
	for k := range t.static {
		seen[k] = struct{}{}
	}
	for _, tags := range t.prefix {
		for k := range tags {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
