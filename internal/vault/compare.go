package vault

import (
	"context"
	"fmt"
)

// VersionPair holds two secret versions for comparison.
type VersionPair struct {
	Path    string
	Mount   string
	VersionA int
	VersionB int
	SecretA  *SecretVersion
	SecretB  *SecretVersion
}

// FetchVersionPair retrieves two versions of a secret from Vault
// and returns them together for diffing.
func FetchVersionPair(ctx context.Context, c *Client, mount, path string, vA, vB int) (*VersionPair, error) {
	secA, err := c.GetSecretVersion(ctx, mount, path, vA)
	if err != nil {
		return nil, fmt.Errorf("fetch version %d: %w", vA, err)
	}

	secB, err := c.GetSecretVersion(ctx, mount, path, vB)
	if err != nil {
		return nil, fmt.Errorf("fetch version %d: %w", vB, err)
	}

	return &VersionPair{
		Path:     path,
		Mount:    mount,
		VersionA: vA,
		VersionB: vB,
		SecretA:  secA,
		SecretB:  secB,
	}, nil
}

// DataA returns the key-value data for version A, never nil.
func (vp *VersionPair) DataA() map[string]string {
	if vp.SecretA == nil {
		return map[string]string{}
	}
	return vp.SecretA.Data
}

// DataB returns the key-value data for version B, never nil.
func (vp *VersionPair) DataB() map[string]string {
	if vp.SecretB == nil {
		return map[string]string{}
	}
	return vp.SecretB.Data
}
