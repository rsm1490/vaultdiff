package vault

import (
	"testing"
)

func TestVersionMeta_Fields(t *testing.T) {
	v := VersionMeta{
		Version:      3,
		CreatedTime:  "2024-01-15T10:00:00Z",
		DeletionTime: "",
		Destroyed:    false,
	}

	if v.Version != 3 {
		t.Errorf("expected Version=3, got %d", v.Version)
	}
	if v.CreatedTime != "2024-01-15T10:00:00Z" {
		t.Errorf("unexpected CreatedTime: %s", v.CreatedTime)
	}
	if v.Destroyed {
		t.Error("expected Destroyed=false")
	}
}

func TestVersionMeta_Destroyed(t *testing.T) {
	v := VersionMeta{
		Version:   1,
		Destroyed: true,
	}
	if !v.Destroyed {
		t.Error("expected Destroyed=true")
	}
}

func TestParseVersionsMap(t *testing.T) {
	raw := map[string]interface{}{
		"1": map[string]interface{}{
			"created_time":  "2024-01-01T00:00:00Z",
			"deletion_time": "",
			"destroyed":     false,
		},
		"2": map[string]interface{}{
			"created_time":  "2024-02-01T00:00:00Z",
			"deletion_time": "2024-03-01T00:00:00Z",
			"destroyed":     true,
		},
	}

	var metas []VersionMeta
	for _, v := range raw {
		entry := v.(map[string]interface{})
		meta := VersionMeta{}
		if ct, ok := entry["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := entry["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if d, ok := entry["destroyed"].(bool); ok {
			meta.Destroyed = d
		}
		metas = append(metas, meta)
	}

	if len(metas) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(metas))
	}

	destroyedCount := 0
	for _, m := range metas {
		if m.Destroyed {
			destroyedCount++
		}
	}
	if destroyedCount != 1 {
		t.Errorf("expected 1 destroyed version, got %d", destroyedCount)
	}
}
