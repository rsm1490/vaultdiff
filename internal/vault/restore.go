package vault

import (
	"context"
	"fmt"
)

// RestoreResult holds the outcome of a secret restore operation.
type RestoreResult struct {
	Path    string
	Version int
	Success bool
	Error   error
}

// Restorer restores a Vault KV v2 secret to a previously deleted version.
type Restorer struct {
	client *Client
}

// NewRestorer creates a new Restorer backed by the given client.
func NewRestorer(c *Client) *Restorer {
	return &Restorer{client: c}
}

// Restore undeletes the specified version of a secret at path under mount.
// It returns a RestoreResult describing whether the operation succeeded.
func (r *Restorer) Restore(ctx context.Context, mount, path string, version int) RestoreResult {
	result := RestoreResult{
		Path:    fmt.Sprintf("%s/%s", mount, path),
		Version: version,
	}

	vaultPath := fmt.Sprintf("%s/undelete/%s", mount, path)
	body := map[string]interface{}{
		"versions": []int{version},
	}

	_, err := r.client.Logical().WriteWithContext(ctx, vaultPath, body)
	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("restore %s version %d: %w", result.Path, version, err)
		return result
	}

	result.Success = true
	return result
}
