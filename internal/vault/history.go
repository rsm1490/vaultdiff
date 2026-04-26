package vault

import (
	"context"
	"fmt"
	"sort"
)

// VersionHistory holds an ordered list of secret versions for a given path.
type VersionHistory struct {
	Path     string
	Versions []VersionMeta
}

// GetHistory retrieves all version metadata for a KV v2 secret path.
func (c *Client) GetHistory(ctx context.Context, mount, secretPath string) (*VersionHistory, error) {
	kvPath := fmt.Sprintf("%s/metadata/%s", mount, secretPath)

	secret, err := c.vault.KVv2(mount).GetMetadata(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("fetching metadata for %s/%s: %w", mount, secretPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no metadata found at %s", kvPath)
	}

	versionsRaw, ok := secret.Raw.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions field in metadata response for %s", kvPath)
	}

	parsed, err := parseVersionsMap(versionsRaw)
	if err != nil {
		return nil, fmt.Errorf("parsing versions for %s: %w", secretPath, err)
	}

	sort.Slice(parsed, func(i, j int) bool {
		return parsed[i].Version < parsed[j].Version
	})

	return &VersionHistory{
		Path:     secretPath,
		Versions: parsed,
	}, nil
}

// Latest returns the highest non-destroyed version, or zero if none exist.
func (h *VersionHistory) Latest() int {
	for i := len(h.Versions) - 1; i >= 0; i-- {
		if !h.Versions[i].Destroyed {
			return h.Versions[i].Version
		}
	}
	return 0
}
