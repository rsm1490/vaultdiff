package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var restoreMount string

var restoreCmd = &cobra.Command{
	Use:   "restore <path> <version>",
	Short: "Restore a deleted secret version in Vault",
	Args:  cobra.ExactArgs(2),
	RunE:  runRestore,
}

func init() {
	restoreCmd.Flags().StringVar(&restoreMount, "mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: must be an integer", args[1])
	}
	if version < 1 {
		return fmt.Errorf("version must be >= 1, got %d", version)
	}

	client := vault.NewClient()
	restorer := vault.NewRestorer(client)

	result := restorer.Restore(context.Background(), restoreMount, path, version)
	if !result.Success {
		return fmt.Errorf("restore failed: %w", result.Error)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully restored %s to version %d\n", result.Path, result.Version)
	return nil
}
