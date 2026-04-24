package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	vc *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
	Mount   string
}

// SecretVersion represents a single version of a KV v2 secret.
type SecretVersion struct {
	Version  int
	Data     map[string]string
	Metadata map[string]interface{}
}

// NewClient creates a new authenticated Vault client.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()
	if cfg.Address != "" {
		vcfg.Address = cfg.Address
	}

	vc, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	if cfg.Token != "" {
		vc.SetToken(cfg.Token)
	}

	return &Client{vc: vc}, nil
}

// GetSecretVersion fetches a specific version of a KV v2 secret.
// If version is 0, the latest version is returned.
func (c *Client) GetSecretVersion(ctx context.Context, mount, path string, version int) (*SecretVersion, error) {
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	secret, err := c.vc.KVv2(mount).GetVersion(ctx, path, version)
	if err != nil {
		return nil, fmt.Errorf("reading secret %s@v%d: %w", path, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found: %s", path)
	}

	data := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		data[k] = fmt.Sprintf("%v", v)
	}

	return &SecretVersion{
		Version:  version,
		Data:     data,
		Metadata: map[string]interface{}{"created_time": secret.VersionMetadata.CreatedTime},
	}, nil
}
