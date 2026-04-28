package vault

import (
	"testing"
)

func TestRollbackResult_Fields(t *testing.T) {
	r := &RollbackResult{
		Path:        "secret/myapp/config",
		FromVersion: 5,
		ToVersion:   3,
		Success:     true,
	}

	if r.Path != "secret/myapp/config" {
		t.Errorf("expected Path %q, got %q", "secret/myapp/config", r.Path)
	}
	if r.FromVersion != 5 {
		t.Errorf("expected FromVersion 5, got %d", r.FromVersion)
	}
	if r.ToVersion != 3 {
		t.Errorf("expected ToVersion 3, got %d", r.ToVersion)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestRollbackResult_FailureState(t *testing.T) {
	r := &RollbackResult{
		Path:        "secret/myapp/db",
		FromVersion: 2,
		ToVersion:   1,
		Success:     false,
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
}

func TestNewRollbacker_NotNil(t *testing.T) {
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("NewClient() error: %v", err)
	}

	rb := NewRollbacker(client)
	if rb == nil {
		t.Fatal("expected non-nil Rollbacker")
	}
}

func TestNewRollbacker_StoresClient(t *testing.T) {
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("NewClient() error: %v", err)
	}

	rb := NewRollbacker(client)
	if rb.client != client {
		t.Error("expected Rollbacker to store the provided client")
	}
}
