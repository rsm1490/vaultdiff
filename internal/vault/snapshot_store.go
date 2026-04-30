package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SnapshotStore persists and loads Snapshot objects from disk.
type SnapshotStore struct {
	Dir string
}

// NewSnapshotStore creates a store rooted at dir.
func NewSnapshotStore(dir string) *SnapshotStore {
	return &SnapshotStore{Dir: dir}
}

// Save writes the snapshot as a JSON file named by ID and timestamp.
func (ss *SnapshotStore) Save(s *Snapshot) (string, error) {
	if err := os.MkdirAll(ss.Dir, 0o755); err != nil {
		return "", fmt.Errorf("snapshot store mkdir: %w", err)
	}
	filename := fmt.Sprintf("%s_%s.json", s.ID, s.TakenAt.Format("20060102T150405Z"))
	path := filepath.Join(ss.Dir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("snapshot store create: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return "", fmt.Errorf("snapshot store encode: %w", err)
	}
	return path, nil
}

// Load reads a snapshot from a file path.
func (ss *SnapshotStore) Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot store open: %w", err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot store decode: %w", err)
	}
	return &s, nil
}

// snapshotFileName is a helper used in tests.
func snapshotFileName(id string, at time.Time) string {
	return fmt.Sprintf("%s_%s.json", id, at.Format("20060102T150405Z"))
}
