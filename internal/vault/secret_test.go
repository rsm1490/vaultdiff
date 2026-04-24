package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestParseSecretVersion_ValidData(t *testing.T) {
	secret := &vaultapi.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"username": "admin",
				"password": "s3cr3t",
			},
			"metadata": map[string]interface{}{
				"version": float64(3),
			},
		},
	}

	sv, err := parseSecretVersion(secret)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if sv.Version != 3 {
		t.Errorf("expected version 3, got %d", sv.Version)
	}
	if sv.Data["username"] != "admin" {
		t.Errorf("expected username 'admin', got %q", sv.Data["username"])
	}
	if sv.Data["password"] != "s3cr3t" {
		t.Errorf("expected password 's3cr3t', got %q", sv.Data["password"])
	}
}

func TestParseSecretVersion_MissingDataField(t *testing.T) {
	secret := &vaultapi.Secret{
		Data: map[string]interface{}{
			"metadata": map[string]interface{}{},
		},
	}

	_, err := parseSecretVersion(secret)
	if err == nil {
		t.Fatal("expected error for missing 'data' field, got nil")
	}
}

func TestParseSecretVersion_NoMetadata(t *testing.T) {
	secret := &vaultapi.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"key": "value",
			},
		},
	}

	sv, err := parseSecretVersion(secret)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if sv.Version != 0 {
		t.Errorf("expected version 0 when metadata absent, got %d", sv.Version)
	}
	if sv.Data["key"] != "value" {
		t.Errorf("expected key 'value', got %q", sv.Data["key"])
	}
}
