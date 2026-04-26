package vault

import (
	"testing"
	"time"
)

func makeHistory(versions []VersionMeta) *VersionHistory {
	return &VersionHistory{
		Path:     "secret/myapp/config",
		Versions: versions,
	}
}

func TestVersionHistory_Latest_ReturnsHighest(t *testing.T) {
	h := makeHistory([]VersionMeta{
		{Version: 1, CreatedAt: time.Now(), Destroyed: false},
		{Version: 2, CreatedAt: time.Now(), Destroyed: false},
		{Version: 3, CreatedAt: time.Now(), Destroyed: false},
	})
	if got := h.Latest(); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}

func TestVersionHistory_Latest_SkipsDestroyed(t *testing.T) {
	h := makeHistory([]VersionMeta{
		{Version: 1, CreatedAt: time.Now(), Destroyed: false},
		{Version: 2, CreatedAt: time.Now(), Destroyed: false},
		{Version: 3, CreatedAt: time.Now(), Destroyed: true},
	})
	if got := h.Latest(); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestVersionHistory_Latest_AllDestroyed(t *testing.T) {
	h := makeHistory([]VersionMeta{
		{Version: 1, Destroyed: true},
		{Version: 2, Destroyed: true},
	})
	if got := h.Latest(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestVersionHistory_Latest_Empty(t *testing.T) {
	h := makeHistory([]VersionMeta{})
	if got := h.Latest(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestVersionHistory_Path(t *testing.T) {
	h := makeHistory(nil)
	if h.Path != "secret/myapp/config" {
		t.Errorf("unexpected path: %s", h.Path)
	}
}
