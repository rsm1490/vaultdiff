package vault

import (
	"testing"
)

func makeSecretVersion(data map[string]string) *SecretVersion {
	return &SecretVersion{
		Data:    data,
		Version: 1,
	}
}

func TestVersionPair_DataA_ReturnsData(t *testing.T) {
	vp := &VersionPair{
		SecretA: makeSecretVersion(map[string]string{"key": "val"}),
	}
	got := vp.DataA()
	if got["key"] != "val" {
		t.Errorf("expected 'val', got %q", got["key"])
	}
}

func TestVersionPair_DataB_ReturnsData(t *testing.T) {
	vp := &VersionPair{
		SecretB: makeSecretVersion(map[string]string{"foo": "bar"}),
	}
	got := vp.DataB()
	if got["foo"] != "bar" {
		t.Errorf("expected 'bar', got %q", got["foo"])
	}
}

func TestVersionPair_DataA_NilSecretReturnsEmpty(t *testing.T) {
	vp := &VersionPair{SecretA: nil}
	got := vp.DataA()
	if got == nil {
		t.Fatal("expected non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestVersionPair_DataB_NilSecretReturnsEmpty(t *testing.T) {
	vp := &VersionPair{SecretB: nil}
	got := vp.DataB()
	if got == nil {
		t.Fatal("expected non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestVersionPair_Fields(t *testing.T) {
	vp := &VersionPair{
		Path:     "myapp/config",
		Mount:    "secret",
		VersionA: 2,
		VersionB: 5,
	}
	if vp.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", vp.Path)
	}
	if vp.Mount != "secret" {
		t.Errorf("unexpected Mount: %s", vp.Mount)
	}
	if vp.VersionA != 2 || vp.VersionB != 5 {
		t.Errorf("unexpected versions: %d, %d", vp.VersionA, vp.VersionB)
	}
}
