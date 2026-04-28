package vault

import (
	"context"
	"fmt"
	"strings"
)

// SearchResult holds a matching secret path and the version it was found in.
type SearchResult struct {
	Path    string
	Version int
	Key     string
	Value   string
}

// SearchOptions configures a secret key/value search.
type SearchOptions struct {
	Mount     string
	KeyPrefix string
	Version   int // 0 means latest
}

// Searcher can search secrets within a Vault KV mount.
type Searcher struct {
	client *Client
}

// NewSearcher creates a Searcher backed by the given client.
func NewSearcher(c *Client) *Searcher {
	return &Searcher{client: c}
}

// SearchPath searches all keys in the secret at path for keys matching the
// given prefix. If opts.Version is 0 the latest version is read.
func (s *Searcher) SearchPath(ctx context.Context, path string, opts SearchOptions) ([]SearchResult, error) {
	mount := opts.Mount
	if mount == "" {
		mount = "secret"
	}

	fullPath := fmt.Sprintf("%s/data/%s", mount, path)
	var secret *SecretVersion
	var err error

	if opts.Version > 0 {
		secret, err = s.client.ReadVersion(ctx, fullPath, opts.Version)
	} else {
		secret, err = s.client.ReadLatest(ctx, fullPath)
	}
	if err != nil {
		return nil, fmt.Errorf("search: read %s: %w", path, err)
	}
	if secret == nil {
		return nil, nil
	}

	var results []SearchResult
	for k, v := range secret.Data {
		if opts.KeyPrefix == "" || strings.HasPrefix(k, opts.KeyPrefix) {
			results = append(results, SearchResult{
				Path:    path,
				Version: secret.Version,
				Key:     k,
				Value:   v,
			})
		}
	}
	return results, nil
}
