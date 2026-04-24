package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddr  string
	vaultToken string
	vaultMount string
)

// rootCmd is the base command for vaultdiff.
var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and audit changes between HashiCorp Vault secret versions",
	Long: `vaultdiff is a CLI tool that lets you compare different versions
of secrets stored in HashiCorp Vault KV v2, making it easy to audit
and review configuration changes over time.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&vaultAddr, "vault-addr", "",
		"Vault server address (defaults to VAULT_ADDR env var)",
	)
	rootCmd.PersistentFlags().StringVar(
		&vaultToken, "vault-token", "",
		"Vault token (defaults to VAULT_TOKEN env var)",
	)
	rootCmd.PersistentFlags().StringVar(
		&vaultMount, "mount", "secret",
		"KV v2 mount path",
	)

	// Allow env var fallbacks
	if vaultAddr == "" {
		vaultAddr = os.Getenv("VAULT_ADDR")
	}
	if vaultToken == "" {
		vaultToken = os.Getenv("VAULT_TOKEN")
	}
}
