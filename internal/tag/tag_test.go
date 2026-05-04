package tag_test

import (
	"testing"

	"github.com/vaultwatch/internal/tag"
)

func newTagger() *tag.Tagger {
	return tag.New(map[string]string{
		"env":  "production",
		"team": "platform",
	})
}

func TestTag_StaticTagsAlwaysPresent(t *testing.T) {
	tr := newTagger()
	tags := tr.Tag("secret/data/foo")
	if tags["env"] != "production" {
		t.Errorf("expected env=production, got %q", tags["env"])
	}
	if tags["team"] != "platform" {
		t.Errorf("expected team=platform, got %q", tags["team"])
	}
}

func TestTag_PrefixTagsAppliedOnMatch(t *testing.T) {
	tr := newTagger()
	tr.AddPrefix("secret/data/db", map[string]string{"service": "database"})

	tags := tr.Tag("secret/data/db/postgres")
	if tags["service"] != "database" {
		t.Errorf("expected service=database, got %q", tags["service"])
	}
}

func TestTag_PrefixTagsNotAppliedOnMismatch(t *testing.T) {
	tr := newTagger()
	tr.AddPrefix("secret/data/db", map[string]string{"service": "database"})

	tags := tr.Tag("secret/data/cache/redis")
	if _, ok := tags["service"]; ok {
		t.Error("expected no service tag for non-matching prefix")
	}
}

func TestTag_PrefixOverridesStatic(t *testing.T) {
	tr := newTagger()
	tr.AddPrefix("secret/data/staging", map[string]string{"env": "staging"})

	tags := tr.Tag("secret/data/staging/api")
	if tags["env"] != "staging" {
		t.Errorf("expected env=staging (prefix override), got %q", tags["env"])
	}
}

func TestTag_DoesNotMutateOriginalMap(t *testing.T) {
	orig := map[string]string{"env": "production"}
	tr := tag.New(orig)

	tags := tr.Tag("any/lease")
	tags["env"] = "mutated"

	if tr.Tag("any/lease")["env"] != "production" {
		t.Error("Tagger internal state was mutated via returned map")
	}
}

func TestKeys_IncludesAllRegisteredKeys(t *testing.T) {
	tr := tag.New(map[string]string{"env": "prod"})
	tr.AddPrefix("secret/", map[string]string{"region": "us-east-1"})

	keys := tr.Keys()
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}
	for _, expected := range []string{"env", "region"} {
		if !keySet[expected] {
			t.Errorf("expected key %q in Keys() result", expected)
		}
	}
}

func TestNew_EmptyStaticTags(t *testing.T) {
	tr := tag.New(nil)
	tags := tr.Tag("secret/data/foo")
	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}
