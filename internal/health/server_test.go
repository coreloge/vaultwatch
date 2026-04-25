package health_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/health"
)

// freePort returns an available TCP port on localhost.
func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	defer l.Close()
	return fmt.Sprintf("127.0.0.1:%d", l.Addr().(*net.TCPAddr).Port)
}

func TestServer_StartsAndResponds(t *testing.T) {
	addr := freePort(t)
	h := health.New(&mockChecker{healthy: true})
	srv := health.NewServer(addr, h)
	srv.Start()
	defer srv.Shutdown()

	// Allow a moment for the listener to be ready.
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s/health", addr))
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_ShutdownClean(t *testing.T) {
	addr := freePort(t)
	h := health.New(&mockChecker{healthy: true})
	srv := health.NewServer(addr, h)
	srv.Start()
	time.Sleep(30 * time.Millisecond)
	// Should not panic or block.
	srv.Shutdown()
}
