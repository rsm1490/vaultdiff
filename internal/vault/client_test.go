package vault

import (
	"testing"
)

func TestNewClient_DefaultAddress(t *testing.T) {
	client, err := NewClient(Config{
		Token: "test-token",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.vc == nil {
		t.Fatal("expected underlying vault client to be initialized")
	}
}

func TestNewClient_CustomAddress(t *testing.T) {
	client, err := NewClient(Config{
		Address: "https://vault.example.com:8200",
		Token:   "s.sometoken",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client.vc.Address() != "https://vault.example.com:8200" {
		t.Errorf("expected custom address, got: %s", client.vc.Address())
	}
}

func TestNewClient_TokenSet(t *testing.T) {
	token := "hvs.testtoken123"
	client, err := NewClient(Config{Token: token})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.vc.Token() != token {
		t.Errorf("expected token %q, got %q", token, client.vc.Token())
	}
}

func TestSecretVersion_DataMapping(t *testing.T) {
	sv := &SecretVersion{
		Version: 3,
		Data: map[string]string{
			"username": "admin",
			"password": "secret",
		},
		Metadata: map[string]interface{}{},
	}

	if sv.Version != 3 {
		t.Errorf("expected version 3, got %d", sv.Version)
	}
	if sv.Data["username"] != "admin" {
		t.Errorf("expected username=admin, got %s", sv.Data["username"])
	}
	if sv.Data["password"] != "secret" {
		t.Errorf("expected password=secret, got %s", sv.Data["password"])
	}
}
