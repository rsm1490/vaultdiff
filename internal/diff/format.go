package diff

import (
	"fmt"
	"io"
	"strings"
)

// OutputFormat defines the output format for diff results.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

// Formatter writes diff results to a writer in a specific format.
type Formatter struct {
	Writer io.Writer
	Format OutputFormat
	Colorize bool
}

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorYellow = "\033[33m"
	colorReset = "\033[0m"
)

// Write outputs the diff results to the formatter's writer.
func (f *Formatter) Write(results []Result) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(results)
	default:
		return f.writeText(results)
	}
}

func (f *Formatter) writeText(results []Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.Writer, "No changes detected.")
		return err
	}

	for _, r := range results {
		var line string
		switch r.Change {
		case ChangeAdded:
			line = fmt.Sprintf("+ %s: %s", r.Key, r.NewValue)
			if f.Colorize {
				line = colorGreen + line + colorReset
			}
		case ChangeRemoved:
			line = fmt.Sprintf("- %s: %s", r.Key, r.OldValue)
			if f.Colorize {
				line = colorRed + line + colorReset
			}
		case ChangeModified:
			line = fmt.Sprintf("~ %s: %s -> %s", r.Key, r.OldValue, r.NewValue)
			if f.Colorize {
				line = colorYellow + line + colorReset
			}
		}
		if _, err := fmt.Fprintln(f.Writer, line); err != nil {
			return err
		}
	}

	summary := Summary(results)
	sep := strings.Repeat("-", 40)
	_, err := fmt.Fprintf(f.Writer, "%s\n%s\n", sep, summary)
	return err
}

func (f *Formatter) writeJSON(results []Result) error {
	type jsonResult struct {
		Key      string `json:"key"`
		Change   string `json:"change"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}

	output := make([]jsonResult, 0, len(results))
	for _, r := range results {
		output = append(output, jsonResult{
			Key:      r.Key,
			Change:   string(r.Change),
			OldValue: r.OldValue,
			NewValue: r.NewValue,
		})
	}

	data, err := marshalJSON(output)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.Writer, string(data))
	return err
}
