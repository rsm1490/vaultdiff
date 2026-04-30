package vault

import (
	"testing"
	"time"
)

func makeSnapshotEntry(path string, version int) SnapshotEntry {
	return SnapshotEntry{
		Path:    path,
		Version: version,
		Data:    map[string]string{"key": "value"},
		TakenAt: time.Now().UTC(),
	}
}

func TestSnapshotEntry_Fields(t *testing.T) {
	e := makeSnapshotEntry("secret/data/foo", 3)
	if e.Path != "secret/data/foo" {
		t.Errorf("expected path secret/data/foo, got %s", e.Path)
	}
	if e.Version != 3 {
		t.Errorf("expected version 3, got %d", e.Version)
	}
	if e.Data["key"] != "value" {
		t.Errorf("expected data key=value")
	}
	if e.TakenAt.IsZero() {
		t.Error("expected TakenAt to be set")
	}
}

func TestNewSnapshot_SetsID(t *testing.T) {
	entries := []SnapshotEntry{makeSnapshotEntry("secret/data/foo", 1)}
	s := NewSnapshot("snap-001", entries)
	if s.ID != "snap-001" {
		t.Errorf("expected ID snap-001, got %s", s.ID)
	}
}

func TestNewSnapshot_SetsEntries(t *testing.T) {
	entries := []SnapshotEntry{
		makeSnapshotEntry("secret/data/a", 1),
		makeSnapshotEntry("secret/data/b", 2),
	}
	s := NewSnapshot("snap-002", entries)
	if len(s.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(s.Entries))
	}
}

func TestNewSnapshot_TakenAtSet(t *testing.T) {
	s := NewSnapshot("snap-003", nil)
	if s.TakenAt.IsZero() {
		t.Error("expected TakenAt to be set")
	}
}

func TestNewSnapshotter_NotNil(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1:8200", "token")
	sn := NewSnapshotter(c)
	if sn == nil {
		t.Error("expected non-nil Snapshotter")
	}
}

func TestNewSnapshotter_StoresClient(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1:8200", "token")
	sn := NewSnapshotter(c)
	if sn.client != c {
		t.Error("expected client to be stored")
	}
}
