package tag_test

import (
	"sync"
	"testing"

	"github.com/vaultwatch/internal/tag"
)

// TestConcurrentTagAndAddPrefix ensures the Tagger is safe for concurrent use.
func TestConcurrentTagAndAddPrefix(t *testing.T) {
	tr := tag.New(map[string]string{"env": "prod"})

	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			tr.AddPrefix("secret/data/", map[string]string{"idx": "val"})
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = tr.Tag("secret/data/foo")
		}()
	}

	wg.Wait()
}

// TestConcurrentKeys ensures Keys() is safe under concurrent mutations.
func TestConcurrentKeys(t *testing.T) {
	tr := tag.New(map[string]string{"a": "1"})

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			tr.AddPrefix("pfx/", map[string]string{"b": "2"})
		}()
		go func() {
			defer wg.Done()
			_ = tr.Keys()
		}()
	}
	wg.Wait()
}
