package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var exportCmd = &cobra.Command{
	Use:   "export <path> <version>",
	Short: "Export a secret version to JSON or env format",
	Args:  cobra.ExactArgs(2),
	RunE:  runExport,
}

func init() {
	exportCmd.Flags().String("mount", "secret", "KV mount path")
	exportCmd.Flags().StringP("format", "f", "json", "Output format: json or env")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	path := args[0]
	versionStr := args[1]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", versionStr, err)
	}

	mount, _ := cmd.Flags().GetString("mount")
	formatStr, _ := cmd.Flags().GetString("format")

	var format vault.ExportFormat
	switch formatStr {
	case "json":
		format = vault.ExportFormatJSON
	case "env":
		format = vault.ExportFormatEnv
	default:
		return fmt.Errorf("unsupported format %q: use json or env", formatStr)
	}

	token := os.Getenv("VAULT_TOKEN")
	addr := os.Getenv("VAULT_ADDR")

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	exporter := vault.NewExporter(client)
	_, err = exporter.Export(path, version, format, mount, os.Stdout)
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	return nil
}
