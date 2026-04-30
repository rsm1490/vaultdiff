package vault

import (
	"fmt"
	"time"
)

// SnapshotEntry represents a captured state of a single secret version.
type SnapshotEntry struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Data    map[string]string `json:"data"`
	TakenAt time.Time         `json:"taken_at"`
}

// Snapshot holds a collection of secret entries captured at a point in time.
type Snapshot struct {
	ID      string          `json:"id"`
	TakenAt time.Time       `json:"taken_at"`
	Entries []SnapshotEntry `json:"entries"`
}

// Snapshotter captures secret state from Vault.
type Snapshotter struct {
	client *Client
}

// NewSnapshotter creates a Snapshotter backed by the given client.
func NewSnapshotter(c *Client) *Snapshotter {
	return &Snapshotter{client: c}
}

// Capture reads the secret at path/version and returns a SnapshotEntry.
func (s *Snapshotter) Capture(mount, path string, version int) (*SnapshotEntry, error) {
	secret, err := s.client.ReadSecretVersion(mount, path, version)
	if err != nil {
		return nil, fmt.Errorf("snapshot capture: %w", err)
	}

	entry := &SnapshotEntry{
		Path:    fmt.Sprintf("%s/%s", mount, path),
		Version: version,
		Data:    secret.Data,
		TakenAt: time.Now().UTC(),
	}
	return entry, nil
}

// NewSnapshot bundles multiple entries into a labelled Snapshot.
func NewSnapshot(id string, entries []SnapshotEntry) *Snapshot {
	return &Snapshot{
		ID:      id,
		TakenAt: time.Now().UTC(),
		Entries: entries,
	}
}
