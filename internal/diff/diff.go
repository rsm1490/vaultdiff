package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType represents the type of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference between two secret versions.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two secret versions.
type Result struct {
	Changes []Change
	HasDiff bool
}

// Compare computes the diff between two maps of secret data.
func Compare(oldData, newData map[string]interface{}) *Result {
	result := &Result{}

	allKeys := mergeKeys(oldData, newData)
	sort.Strings(allKeys)

	for _, key := range allKeys {
		oldVal, oldExists := oldData[key]
		newVal, newExists := newData[key]

		switch {
		case !oldExists && newExists:
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Added,
				NewValue: fmt.Sprintf("%v", newVal),
			})
			result.HasDiff = true
		case oldExists && !newExists:
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: fmt.Sprintf("%v", oldVal),
			})
			result.HasDiff = true
		default:
			oldStr := fmt.Sprintf("%v", oldVal)
			newStr := fmt.Sprintf("%v", newVal)
			if oldStr != newStr {
				result.Changes = append(result.Changes, Change{
					Key:      key,
					Type:     Modified,
					OldValue: oldStr,
					NewValue: newStr,
				})
				result.HasDiff = true
			} else {
				result.Changes = append(result.Changes, Change{
					Key:  key,
					Type: Unchanged,
				})
			}
		}
	}

	return result
}

// mergeKeys returns the union of keys from both maps.
func mergeKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}

// Summary returns a human-readable summary of the diff result.
func (r *Result) Summary() string {
	if !r.HasDiff {
		return "No differences found."
	}
	var sb strings.Builder
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(&sb, "+ %s: %q\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(&sb, "- %s: %q\n", c.Key, c.OldValue)
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %q -> %q\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return sb.String()
}
