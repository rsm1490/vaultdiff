package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// SecretVersion represents a specific version of a KV v2 secret.
type SecretVersion struct {
	Version  int
	Data     map[string]string
	Metadata map[string]interface{}
}

// GetSecretVersion retrieves a specific version of a KV v2 secret at the given path.
// If version is 0, the latest version is returned.
func (c *Client) GetSecretVersion(ctx context.Context, mount, path string, version int) (*SecretVersion, error) {
	kvPath := fmt.Sprintf("%s/data/%s", mount, path)

	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	secret, err := c.vault.Logical().ReadWithDataWithContext(ctx, kvPath, params)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q (version %d): %w", path, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %q not found", path)
	}

	return parseSecretVersion(secret)
}

func parseSecretVersion(secret *vaultapi.Secret) (*SecretVersion, error) {
	rawData, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("secret response missing 'data' field")
	}

	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected 'data' field type")
	}

	strData := make(map[string]string, len(dataMap))
	for k, v := range dataMap {
		strData[k] = fmt.Sprintf("%v", v)
	}

	meta, _ := secret.Data["metadata"].(map[string]interface{})

	var version int
	if meta != nil {
		if v, ok := meta["version"].(float64); ok {
			version = int(v)
		}
	}

	return &SecretVersion{
		Version:  version,
		Data:     strData,
		Metadata: meta,
	}, nil
}
