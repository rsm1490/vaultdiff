package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Path        string
	FromVersion int
	ToVersion   int
	Success     bool
}

// Rollbacker performs version rollback operations on KV v2 secrets.
type Rollbacker struct {
	client *vaultapi.Client
}

// NewRollbacker creates a Rollbacker using the provided Vault client.
func NewRollbacker(client *vaultapi.Client) *Rollbacker {
	return &Rollbacker{client: client}
}

// Rollback copies the data from targetVersion of the secret at path and writes
// it as a new version, effectively rolling back to that state.
func (r *Rollbacker) Rollback(ctx context.Context, mount, path string, targetVersion int) (*RollbackResult, error) {
	kvPath := fmt.Sprintf("%s/data/%s", mount, path)
	metaPath := fmt.Sprintf("%s/metadata/%s", mount, path)

	// Fetch current metadata to determine the live version.
	meta, err := r.client.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %s: %w", path, err)
	}
	if meta == nil {
		return nil, fmt.Errorf("secret not found: %s", path)
	}

	currentVersion, _ := meta.Data["current_version"].(int)

	// Read the target version.
	secret, err := r.client.Logical().ReadWithDataWithContext(ctx, kvPath, map[string][]string{
		"version": {fmt.Sprintf("%d", targetVersion)},
	})
	if err != nil {
		return nil, fmt.Errorf("reading version %d of %s: %w", targetVersion, path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("version %d not found for secret: %s", targetVersion, path)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format for version %d of %s", targetVersion, path)
	}

	// Write the historical data as a new version.
	_, err = r.client.Logical().WriteWithContext(ctx, kvPath, map[string]interface{}{
		"data": data,
	})
	if err != nil {
		return nil, fmt.Errorf("writing rollback data to %s: %w", path, err)
	}

	return &RollbackResult{
		Path:        path,
		FromVersion: currentVersion,
		ToVersion:   targetVersion,
		Success:     true,
	}, nil
}
