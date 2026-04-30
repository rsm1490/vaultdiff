package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatEnv  ExportFormat = "env"
)

// ExportResult holds the result of a secret export operation.
type ExportResult struct {
	Path      string
	Version   int
	Data      map[string]string
	ExportedAt time.Time
	Format    ExportFormat
}

// Exporter handles exporting secret versions to various formats.
type Exporter struct {
	client *Client
}

// NewExporter creates a new Exporter backed by the given client.
func NewExporter(client *Client) *Exporter {
	return &Exporter{client: client}
}

// Export fetches the secret at path/version and writes it to w in the given format.
func (e *Exporter) Export(path string, version int, format ExportFormat, mount string, w io.Writer) (*ExportResult, error) {
	sv, err := e.client.GetSecretVersion(path, version, mount)
	if err != nil {
		return nil, fmt.Errorf("export: fetch secret %q version %d: %w", path, version, err)
	}

	result := &ExportResult{
		Path:       path,
		Version:    version,
		Data:       sv.Data,
		ExportedAt: time.Now().UTC(),
		Format:     format,
	}

	switch format {
	case ExportFormatJSON:
		if err := writeExportJSON(result, w); err != nil {
			return nil, err
		}
	case ExportFormatEnv:
		writeExportEnv(result, w)
	default:
		return nil, fmt.Errorf("export: unsupported format %q", format)
	}

	return result, nil
}

func writeExportJSON(r *ExportResult, w io.Writer) error {
	payload := map[string]interface{}{
		"path":        r.Path,
		"version":     r.Version,
		"exported_at": r.ExportedAt.Format(time.RFC3339),
		"data":        r.Data,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func writeExportEnv(r *ExportResult, w io.Writer) {
	for k, v := range r.Data {
		fmt.Fprintf(w, "%s=%q\n", k, v)
	}
}
