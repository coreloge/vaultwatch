// Package labelset provides structured key-value label management for lease
// events, enabling consistent tagging across alert payloads and audit logs.
package labelset

import (
	"fmt"
	"sort"
	"strings"
)

// LabelSet holds an immutable set of string key-value labels.
type LabelSet struct {
	labels map[string]string
}

// New creates a LabelSet from the provided key-value pairs.
// Panics if an odd number of arguments is supplied.
func New(kvs ...string) LabelSet {
	if len(kvs)%2 != 0 {
		panic(fmt.Sprintf("labelset.New: odd number of arguments (%d)", len(kvs)))
	}
	m := make(map[string]string, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return LabelSet{labels: m}
}

// FromMap creates a LabelSet from an existing map. The map is copied.
func FromMap(m map[string]string) LabelSet {
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return LabelSet{labels: copy}
}

// Get returns the value for key and whether it was present.
func (ls LabelSet) Get(key string) (string, bool) {
	v, ok := ls.labels[key]
	return v, ok
}

// Merge returns a new LabelSet combining ls with other.
// Keys in other override keys in ls.
func (ls LabelSet) Merge(other LabelSet) LabelSet {
	merged := make(map[string]string, len(ls.labels)+len(other.labels))
	for k, v := range ls.labels {
		merged[k] = v
	}
	for k, v := range other.labels {
		merged[k] = v
	}
	return LabelSet{labels: merged}
}

// ToMap returns a shallow copy of the underlying label map.
func (ls LabelSet) ToMap() map[string]string {
	copy := make(map[string]string, len(ls.labels))
	for k, v := range ls.labels {
		copy[k] = v
	}
	return copy
}

// Len returns the number of labels.
func (ls LabelSet) Len() int { return len(ls.labels) }

// String returns a deterministic comma-separated key=value representation.
func (ls LabelSet) String() string {
	keys := make([]string, 0, len(ls.labels))
	for k := range ls.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+ls.labels[k])
	}
	return strings.Join(parts, ",")
}
