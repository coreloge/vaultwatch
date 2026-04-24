package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// LeaseInfo holds metadata about a Vault secret lease.
type LeaseInfo struct {
	LeaseID   string
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
}

// Client wraps the Vault API client with lease-inspection helpers.
type Client struct {
	api *vaultapi.Client
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	c.SetToken(token)

	return &Client{api: c}, nil
}

// LookupLease returns lease information for the given lease ID.
func (c *Client) LookupLease(ctx context.Context, leaseID string) (*LeaseInfo, error) {
	secret, err := c.api.Sys().LookupLeaseWithContext(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("looking up lease %q: %w", leaseID, err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data returned for lease %q", leaseID)
	}

	ttlRaw, ok := secret.Data["ttl"]
	if !ok {
		return nil, fmt.Errorf("ttl missing from lease data for %q", leaseID)
	}

	ttlSeconds, ok := ttlRaw.(float64)
	if !ok {
		return nil, fmt.Errorf("unexpected ttl type for lease %q", leaseID)
	}

	ttl := time.Duration(ttlSeconds) * time.Second
	expiresAt := time.Now().Add(ttl)

	pathRaw, _ := secret.Data["path"]
	path, _ := pathRaw.(string)

	return &LeaseInfo{
		LeaseID:   leaseID,
		Path:      path,
		ExpiresAt: expiresAt,
		TTL:       ttl,
	}, nil
}

// IsHealthy checks whether the Vault server is reachable and unsealed.
func (c *Client) IsHealthy(ctx context.Context) error {
	health, err := c.api.Sys().HealthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("vault health check failed: %w", err)
	}
	if health.Sealed {
		return fmt.Errorf("vault is sealed")
	}
	return nil
}
