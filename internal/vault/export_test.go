package vault

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestExportResult_Fields(t *testing.T) {
	r := &ExportResult{
		Path:    "secret/myapp",
		Version: 3,
		Data:    map[string]string{"KEY": "val"},
		Format:  ExportFormatJSON,
	}
	if r.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("expected version 3, got %d", r.Version)
	}
	if r.Data["KEY"] != "val" {
		t.Errorf("expected KEY=val")
	}
}

func TestNewExporter_NotNil(t *testing.T) {
	c, _ := NewClient("", "")
	e := NewExporter(c)
	if e == nil {
		t.Fatal("expected non-nil Exporter")
	}
}

func TestNewExporter_StoresClient(t *testing.T) {
	c, _ := NewClient("", "")
	e := NewExporter(c)
	if e.client != c {
		t.Error("expected exporter to store client")
	}
}

func TestWriteExportJSON_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	r := &ExportResult{
		Path:    "secret/app",
		Version: 2,
		Data:    map[string]string{"DB_PASS": "s3cr3t"},
		Format:  ExportFormatJSON,
	}
	if err := writeExportJSON(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["path"] != "secret/app" {
		t.Errorf("expected path secret/app, got %v", out["path"])
	}
	if out["version"].(float64) != 2 {
		t.Errorf("expected version 2")
	}
}

func TestWriteExportEnv_KeyValueFormat(t *testing.T) {
	var buf bytes.Buffer
	r := &ExportResult{
		Data: map[string]string{"API_KEY": "abc123"},
	}
	writeExportEnv(r, &buf)
	output := buf.String()
	if !strings.Contains(output, "API_KEY=") {
		t.Errorf("expected API_KEY in env output, got: %s", output)
	}
}

func TestExportFormat_Constants(t *testing.T) {
	if ExportFormatJSON != "json" {
		t.Errorf("expected json, got %s", ExportFormatJSON)
	}
	if ExportFormatEnv != "env" {
		t.Errorf("expected env, got %s", ExportFormatEnv)
	}
}
