package vault

import (
	"testing"
	"time"
)

func TestWatchOptions_Defaults(t *testing.T) {
	opts := DefaultWatchOptions()
	if opts.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", opts.Mount)
	}
	if opts.Interval != 30*time.Second {
		t.Errorf("expected interval 30s, got %v", opts.Interval)
	}
}

func TestNewWatcher_NotNil(t *testing.T) {
	client := &Client{}
	w := NewWatcher(client, DefaultWatchOptions())
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

func TestNewWatcher_StoresClient(t *testing.T) {
	client := &Client{}
	w := NewWatcher(client, DefaultWatchOptions())
	if w.client != client {
		t.Error("expected watcher to store the provided client")
	}
}

func TestNewWatcher_StoresOptions(t *testing.T) {
	client := &Client{}
	opts := WatchOptions{Mount: "kv", Interval: 10 * time.Second}
	w := NewWatcher(client, opts)
	if w.options.Mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", w.options.Mount)
	}
	if w.options.Interval != 10*time.Second {
		t.Errorf("expected interval 10s, got %v", w.options.Interval)
	}
}

func TestWatchEvent_Fields(t *testing.T) {
	now := time.Now().UTC()
	ev := WatchEvent{
		Path:       "myapp/config",
		OldVersion: 2,
		NewVersion: 3,
		ChangedAt:  now,
	}
	if ev.Path != "myapp/config" {
		t.Errorf("unexpected Path: %q", ev.Path)
	}
	if ev.OldVersion != 2 {
		t.Errorf("unexpected OldVersion: %d", ev.OldVersion)
	}
	if ev.NewVersion != 3 {
		t.Errorf("unexpected NewVersion: %d", ev.NewVersion)
	}
	if !ev.ChangedAt.Equal(now) {
		t.Errorf("unexpected ChangedAt: %v", ev.ChangedAt)
	}
}
