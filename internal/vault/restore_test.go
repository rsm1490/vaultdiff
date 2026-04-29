package vault

import (
	"testing"
)

func TestRestoreResult_Fields(t *testing.T) {
	r := RestoreResult{
		Path:    "secret/myapp",
		Version: 3,
		Success: true,
		Error:   nil,
	}

	if r.Path != "secret/myapp" {
		t.Errorf("expected path 'secret/myapp', got %q", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("expected version 3, got %d", r.Version)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.Error != nil {
		t.Errorf("expected nil error, got %v", r.Error)
	}
}

func TestRestoreResult_FailureState(t *testing.T) {
	err := fmt.Errorf("permission denied")
	r := RestoreResult{
		Path:    "secret/myapp",
		Version: 2,
		Success: false,
		Error:   err,
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Error == nil {
		t.Error("expected non-nil error")
	}
}

func TestNewRestorer_NotNil(t *testing.T) {
	c := NewClient()
	restorer := NewRestorer(c)
	if restorer == nil {
		t.Fatal("expected non-nil Restorer")
	}
}

func TestNewRestorer_StoresClient(t *testing.T) {
	c := NewClient()
	restorer := NewRestorer(c)
	if restorer.client != c {
		t.Error("expected restorer to store the provided client")
	}
}
