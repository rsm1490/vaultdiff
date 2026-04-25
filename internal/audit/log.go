package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Entry represents a single audit log entry for a diff operation.
type Entry struct {
	Timestamp  time.Time        `json:"timestamp"`
	Path       string           `json:"path"`
	VersionA   int              `json:"version_a"`
	VersionB   int              `json:"version_b"`
	Changes    []diff.Change    `json:"changes"`
	Summary    diff.Summary     `json:"summary"`
}

// Logger writes audit entries to an output destination.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to w. Pass os.Stdout for console output.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// NewFileLogger opens (or creates) the file at path and returns a Logger
// backed by that file. The caller is responsible for closing the file.
func NewFileLogger(path string) (*Logger, *os.File, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return NewLogger(f), f, nil
}

// Record encodes entry as a JSON line and writes it to the logger's writer.
func (l *Logger) Record(entry Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
