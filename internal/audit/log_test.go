package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultdiff/internal/audit"
	"github.com/yourusername/vaultdiff/internal/diff"
)

func makeEntry() audit.Entry {
	return audit.Entry{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Path:      "secret/data/myapp",
		VersionA:  1,
		VersionB:  2,
		Changes: []diff.Change{
			{Key: "DB_PASS", Type: diff.Modified, OldValue: "old", NewValue: "new"},
		},
		Summary: diff.Summary{Added: 0, Removed: 0, Modified: 1},
	}
}

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	entry := makeEntry()
	if err := logger.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(line, "{") {
		t.Errorf("expected JSON line, got: %s", line)
	}

	var decoded audit.Entry
	if err := json.Unmarshal([]byte(line), &decoded); err != nil {
		t.Fatalf("failed to decode output: %v", err)
	}
	if decoded.Path != entry.Path {
		t.Errorf("path mismatch: got %s, want %s", decoded.Path, entry.Path)
	}
	if decoded.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", decoded.Summary.Modified)
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	entry := makeEntry()
	entry.Timestamp = time.Time{} // zero value

	if err := logger.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &decoded); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if decoded.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestNewLogger_NilWriterDefaultsToStdout(t *testing.T) {
	logger := audit.NewLogger(nil)
	if logger == nil {
		t.Error("expected non-nil logger")
	}
}
