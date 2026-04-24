package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
  token: "s.testtoken"
poll_interval: 60s
warn_threshold: 48h
webhooks:
  - name: slack
    url: "https://hooks.slack.com/test"
    retries: 3
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("expected vault address, got %q", cfg.Vault.Address)
	}
	if cfg.PollInterval != 60*time.Second {
		t.Errorf("expected 60s poll interval, got %v", cfg.PollInterval)
	}
	if cfg.WarnThreshold != 48*time.Hour {
		t.Errorf("expected 48h warn threshold, got %v", cfg.WarnThreshold)
	}
	if len(cfg.Webhooks) != 1 || cfg.Webhooks[0].Name != "slack" {
		t.Errorf("unexpected webhooks: %+v", cfg.Webhooks)
	}
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
  token: "s.tok"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 30*time.Second {
		t.Errorf("expected default poll interval 30s, got %v", cfg.PollInterval)
	}
	if cfg.WarnThreshold != 24*time.Hour {
		t.Errorf("expected default warn threshold 24h, got %v", cfg.WarnThreshold)
	}
}

func TestLoad_MissingVaultAddress(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  token: "s.tok"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}

func TestLoad_MissingToken(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when no auth credentials provided")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
