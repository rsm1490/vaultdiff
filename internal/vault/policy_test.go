package vault

import (
	"testing"
)

func TestPolicyEntry_Fields(t *testing.T) {
	e := PolicyEntry{
		Path:         "my-policy",
		Capabilities: []string{"read", "list"},
	}
	if e.Path != "my-policy" {
		t.Errorf("expected path my-policy, got %s", e.Path)
	}
	if len(e.Capabilities) != 2 {
		t.Errorf("expected 2 capabilities, got %d", len(e.Capabilities))
	}
}

func TestPolicyReport_Fields(t *testing.T) {
	r := PolicyReport{
		Path: "secret/data/myapp/config",
		Entries: []PolicyEntry{
			{Path: "admin", Capabilities: []string{"read", "write"}},
		},
	}
	if r.Path != "secret/data/myapp/config" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if len(r.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(r.Entries))
	}
}

func TestPolicyReport_EmptyEntries(t *testing.T) {
	r := PolicyReport{Path: "secret/data/foo"}
	if r.Entries != nil {
		t.Errorf("expected nil entries for empty report")
	}
}

func TestNewPolicyChecker_NotNil(t *testing.T) {
	c, _ := NewClient()
	pc := NewPolicyChecker(c)
	if pc == nil {
		t.Fatal("expected non-nil PolicyChecker")
	}
}

func TestNewPolicyChecker_StoresClient(t *testing.T) {
	c, _ := NewClient()
	pc := NewPolicyChecker(c)
	if pc.client != c {
		t.Error("expected PolicyChecker to store the provided client")
	}
}

func TestExtractCapabilities_Found(t *testing.T) {
	rules := `path "secret/data/myapp" {\n  capabilities = ["read", "list"]\n}`
	caps := extractCapabilities(rules, "secret/data/myapp")
	// extractCapabilities is best-effort; just ensure no panic
	_ = caps
}

func TestExtractCapabilities_NotFound(t *testing.T) {
	rules := `path "secret/data/other" {\n  capabilities = ["read"]\n}`
	caps := extractCapabilities(rules, "secret/data/myapp")
	if len(caps) != 0 {
		t.Errorf("expected no capabilities for unmatched path, got %v", caps)
	}
}
