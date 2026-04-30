package labelset_test

import (
	"testing"

	"github.com/youorg/vaultwatch/internal/labelset"
)

func newLabels(kvs ...string) labelset.LabelSet {
	return labelset.New(kvs...)
}

func TestNew_StoresKeyValues(t *testing.T) {
	ls := newLabels("env", "prod", "team", "platform")
	if v, ok := ls.Get("env"); !ok || v != "prod" {
		t.Errorf("expected env=prod, got %q ok=%v", v, ok)
	}
	if v, ok := ls.Get("team"); !ok || v != "platform" {
		t.Errorf("expected team=platform, got %q ok=%v", v, ok)
	}
}

func TestNew_OddArgsPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for odd number of args")
		}
	}()
	labelset.New("only-one")
}

func TestGet_MissingKey(t *testing.T) {
	ls := newLabels("a", "1")
	_, ok := ls.Get("missing")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestFromMap_CopiesMap(t *testing.T) {
	m := map[string]string{"x": "1"}
	ls := labelset.FromMap(m)
	m["x"] = "mutated"
	if v, _ := ls.Get("x"); v != "1" {
		t.Errorf("expected original value '1', got %q", v)
	}
}

func TestMerge_OtherOverrides(t *testing.T) {
	a := newLabels("env", "staging", "region", "us-east")
	b := newLabels("env", "prod")
	merged := a.Merge(b)
	if v, _ := merged.Get("env"); v != "prod" {
		t.Errorf("expected env=prod after merge, got %q", v)
	}
	if v, _ := merged.Get("region"); v != "us-east" {
		t.Errorf("expected region=us-east preserved, got %q", v)
	}
}

func TestMerge_DoesNotMutateOriginal(t *testing.T) {
	a := newLabels("env", "staging")
	b := newLabels("env", "prod")
	a.Merge(b)
	if v, _ := a.Get("env"); v != "staging" {
		t.Errorf("original label set mutated: got %q", v)
	}
}

func TestToMap_ReturnsCopy(t *testing.T) {
	ls := newLabels("k", "v")
	m := ls.ToMap()
	m["k"] = "changed"
	if v, _ := ls.Get("k"); v != "v" {
		t.Errorf("LabelSet mutated via ToMap result: got %q", v)
	}
}

func TestLen_ReturnsCount(t *testing.T) {
	ls := newLabels("a", "1", "b", "2", "c", "3")
	if ls.Len() != 3 {
		t.Errorf("expected Len=3, got %d", ls.Len())
	}
}

func TestString_IsDeterministic(t *testing.T) {
	ls := newLabels("z", "last", "a", "first", "m", "mid")
	want := "a=first,m=mid,z=last"
	if got := ls.String(); got != want {
		t.Errorf("String()=%q, want %q", got, want)
	}
}
