package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full vaultwatch daemon configuration.
type Config struct {
	Vault    VaultConfig    `yaml:"vault"`
	Webhooks []WebhookConfig `yaml:"webhooks"`
	PollInterval time.Duration `yaml:"poll_interval"`
	WarnThreshold time.Duration `yaml:"warn_threshold"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	RoleID  string `yaml:"role_id"`
	SecretID string `yaml:"secret_id"`
}

// WebhookConfig defines a single webhook target.
type WebhookConfig struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Retries int               `yaml:"retries"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := &Config{
		PollInterval:  30 * time.Second,
		WarnThreshold: 24 * time.Hour,
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.Token == "" && (c.Vault.RoleID == "" || c.Vault.SecretID == "") {
		return fmt.Errorf("vault.token or vault.role_id + vault.secret_id is required")
	}
	if c.PollInterval <= 0 {
		return fmt.Errorf("poll_interval must be positive")
	}
	if c.WarnThreshold <= 0 {
		return fmt.Errorf("warn_threshold must be positive")
	}
	for i, wh := range c.Webhooks {
		if wh.URL == "" {
			return fmt.Errorf("webhooks[%d].url is required", i)
		}
	}
	return nil
}
