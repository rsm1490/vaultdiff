package vault

import (
	"context"
	"fmt"
	"sort"
)

// VersionMeta holds metadata for a single secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns metadata for all versions of a KVv2 secret.
func (c *Client) ListVersions(ctx context.Context, mount, secretPath string) ([]VersionMeta, error) {
	path := fmt.Sprintf("%s/metadata/%s", mount, secretPath)

	secret, err := c.vault.KVv2(mount).GetMetadata(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("listing versions at %s: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no metadata found at %s", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("metadata response missing 'versions' field")
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for 'versions' field")
	}

	var metas []VersionMeta
	for _, v := range versionsMap {
		entry, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		meta := VersionMeta{}
		if ct, ok := entry["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := entry["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if d, ok := entry["destroyed"].(bool); ok {
			meta.Destroyed = d
		}
		metas = append(metas, meta)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})

	return metas, nil
}
