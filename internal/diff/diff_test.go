package diff

import (
	"strings"
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	old := map[string]interface{}{"key1": "value1", "key2": "value2"}
	new := map[string]interface{}{"key1": "value1", "key2": "value2"}

	result := Compare(old, new)

	if result.HasDiff {
		t.Error("expected no diff, but HasDiff is true")
	}
	if len(result.Changes) != 2 {
		t.Errorf("expected 2 unchanged entries, got %d", len(result.Changes))
	}
	for _, c := range result.Changes {
		if c.Type != Unchanged {
			t.Errorf("expected Unchanged for key %s, got %s", c.Key, c.Type)
		}
	}
}

func TestCompare_AddedKey(t *testing.T) {
	old := map[string]interface{}{"key1": "value1"}
	new := map[string]interface{}{"key1": "value1", "key2": "newvalue"}

	result := Compare(old, new)

	if !result.HasDiff {
		t.Error("expected diff, but HasDiff is false")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "key2" && c.Type == Added && c.NewValue == "newvalue" {
			found = true
		}
	}
	if !found {
		t.Error("expected Added change for key2")
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	old := map[string]interface{}{"key1": "value1", "key2": "oldvalue"}
	new := map[string]interface{}{"key1": "value1"}

	result := Compare(old, new)

	if !result.HasDiff {
		t.Error("expected diff, but HasDiff is false")
	}
	for _, c := range result.Changes {
		if c.Key == "key2" {
			if c.Type != Removed {
				t.Errorf("expected Removed for key2, got %s", c.Type)
			}
			if c.OldValue != "oldvalue" {
				t.Errorf("expected OldValue 'oldvalue', got %s", c.OldValue)
			}
		}
	}
}

func TestCompare_ModifiedKey(t *testing.T) {
	old := map[string]interface{}{"key1": "original"}
	new := map[string]interface{}{"key1": "updated"}

	result := Compare(old, new)

	if !result.HasDiff {
		t.Error("expected diff, but HasDiff is false")
	}
	if result.Changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", result.Changes[0].Type)
	}
}

func TestSummary_NoDiff(t *testing.T) {
	r := &Result{HasDiff: false}
	if r.Summary() != "No differences found." {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestSummary_WithChanges(t *testing.T) {
	old := map[string]interface{}{"a": "1", "b": "old"}
	new := map[string]interface{}{"b": "new", "c": "3"}

	result := Compare(old, new)
	summary := result.Summary()

	if !strings.Contains(summary, "- a") {
		t.Error("expected removed key 'a' in summary")
	}
	if !strings.Contains(summary, "~ b") {
		t.Error("expected modified key 'b' in summary")
	}
	if !strings.Contains(summary, "+ c") {
		t.Error("expected added key 'c' in summary")
	}
}
