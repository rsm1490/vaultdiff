package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var policyCmd = &cobra.Command{
	Use:   "policy <secret-path>",
	Short: "Show which Vault policies cover a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runPolicy,
}

func init() {
	policyCmd.Flags().String("mount", "secret", "KV mount path")
	policyCmd.Flags().String("format", "text", "Output format: text or json")
	rootCmd.AddCommand(policyCmd)
}

func runPolicy(cmd *cobra.Command, args []string) error {
	secretPath := args[0]
	mount, _ := cmd.Flags().GetString("mount")
	format, _ := cmd.Flags().GetString("format")

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	checker := vault.NewPolicyChecker(client)
	report, err := checker.CheckPath(cmd.Context(), secretPath, mount)
	if err != nil {
		return fmt.Errorf("checking policy for %q: %w", secretPath, err)
	}

	switch strings.ToLower(format) {
	case "json":
		return printPolicyJSON(report)
	default:
		return printPolicyText(report)
	}
}

func printPolicyText(r *vault.PolicyReport) error {
	if len(r.Entries) == 0 {
		fmt.Printf("No policies found covering path: %s\n", r.Path)
		return nil
	}
	fmt.Printf("Policies covering path: %s\n\n", r.Path)
	for _, e := range r.Entries {
		fmt.Printf("  Policy : %s\n", e.Path)
		fmt.Printf("  Caps   : %s\n\n", strings.Join(e.Capabilities, ", "))
	}
	return nil
}

func printPolicyJSON(r *vault.PolicyReport) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
