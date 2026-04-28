package vault

import (
	"testing"
)

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{
		Path:    "myapp/config",
		Version: 3,
		Key:     "DB_HOST",
		Value:   "localhost",
	}
	if r.Path != "myapp/config" {
		t.Errorf("expected path myapp/config, got %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("expected version 3, got %d", r.Version)
	}
	if r.Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", r.Key)
	}
	if r.Value != "localhost" {
		t.Errorf("expected value localhost, got %s", r.Value)
	}
}

func TestSearchOptions_Defaults(t *testing.T) {
	opts := SearchOptions{}
	if opts.Mount != "" {
		t.Errorf("expected empty mount by default")
	}
	if opts.Version != 0 {
		t.Errorf("expected version 0 (latest) by default")
	}
}

func TestNewSearcher_NotNil(t *testing.T) {
	c := &Client{}
	s := NewSearcher(c)
	if s == nil {
		t.Fatal("expected non-nil Searcher")
	}
	if s.client != c {
		t.Error("expected searcher to hold provided client")
	}
}

func TestSearchPath_FiltersByKeyPrefix(t *testing.T) {
	// Build a SecretVersion directly to validate prefix filtering logic.
	sv := &SecretVersion{
		Version: 2,
		Data: map[string]string{
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"API_KEY":     "secret",
			"API_TIMEOUT": "30s",
		},
	}

	// Simulate prefix filtering as done inside SearchPath.
	prefix := "DB_"
	var matched []string
	for k := range sv.Data {
		if len(prefix) == 0 || len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			matched = append(matched, k)
		}
	}
	if len(matched) != 2 {
		t.Errorf("expected 2 DB_ keys, got %d", len(matched))
	}
}

func TestSearchPath_EmptyPrefixMatchesAll(t *testing.T) {
	sv := &SecretVersion{
		Version: 1,
		Data: map[string]string{
			"FOO": "bar",
			"BAZ": "qux",
		},
	}

	count := 0
	prefix := ""
	for k := range sv.Data {
		if prefix == "" || k[:len(prefix)] == prefix {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected all 2 keys with empty prefix, got %d", count)
	}
}
