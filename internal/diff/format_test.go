package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatter_WriteText_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatText}

	if err := f.Write([]Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected 'No changes detected', got: %q", buf.String())
	}
}

func TestFormatter_WriteText_Added(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatText}

	results := []Result{
		{Key: "foo", Change: ChangeAdded, NewValue: "bar"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+ foo: bar") {
		t.Errorf("expected added line, got: %q", out)
	}
}

func TestFormatter_WriteText_Removed(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatText}

	results := []Result{
		{Key: "secret", Change: ChangeRemoved, OldValue: "old"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "- secret: old") {
		t.Errorf("expected removed line, got: %q", out)
	}
}

func TestFormatter_WriteText_Modified(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatText}

	results := []Result{
		{Key: "token", Change: ChangeModified, OldValue: "v1", NewValue: "v2"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "~ token: v1 -> v2") {
		t.Errorf("expected modified line, got: %q", out)
	}
}

func TestFormatter_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatJSON}

	results := []Result{
		{Key: "api_key", Change: ChangeAdded, NewValue: "abc123"},
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed []map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected 1 result, got %d", len(parsed))
	}
	if parsed[0]["key"] != "api_key" {
		t.Errorf("expected key 'api_key', got %q", parsed[0]["key"])
	}
	if parsed[0]["change"] != string(ChangeAdded) {
		t.Errorf("expected change 'added', got %q", parsed[0]["change"])
	}
}
