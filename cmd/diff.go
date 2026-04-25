package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"vaultdiff/internal/diff"
	"vaultdiff/internal/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff <secret-path> <version-a> <version-b>",
	Short: "Diff two versions of a Vault KV secret",
	Long: `Compare two versions of a HashiCorp Vault KV v2 secret and display
the differences between their data fields.`,
	Args: cobra.ExactArgs(3),
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolP("show-unchanged", "u", false, "Include unchanged keys in output")
}

func runDiff(cmd *cobra.Command, args []string) error {
	secretPath := args[0]

	versionA, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version-a %q: %w", args[1], err)
	}

	versionB, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid version-b %q: %w", args[2], err)
	}

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	secretA, err := vault.GetSecretVersion(client, secretPath, versionA)
	if err != nil {
		return fmt.Errorf("failed to fetch version %d: %w", versionA, err)
	}

	secretB, err := vault.GetSecretVersion(client, secretPath, versionB)
	if err != nil {
		return fmt.Errorf("failed to fetch version %d: %w", versionB, err)
	}

	result := diff.Compare(secretA.Data, secretB.Data)

	showUnchanged, _ := cmd.Flags().GetBool("show-unchanged")

	fmt.Fprintf(os.Stdout, "Diff for %s (v%d -> v%d):\n\n", secretPath, versionA, versionB)

	if !result.HasDiff {
		fmt.Fprintln(os.Stdout, "No differences found.")
		return nil
	}

	for _, c := range result.Changes {
		switch c.Type {
		case diff.Added:
			fmt.Fprintf(os.Stdout, "+ %s: %q\n", c.Key, c.NewValue)
		case diff.Removed:
			fmt.Fprintf(os.Stdout, "- %s: %q\n", c.Key, c.OldValue)
		case diff.Modified:
			fmt.Fprintf(os.Stdout, "~ %s: %q -> %q\n", c.Key, c.OldValue, c.NewValue)
		case diff.Unchanged:
			if showUnchanged {
				fmt.Fprintf(os.Stdout, "  %s (unchanged)\n", c.Key)
			}
		}
	}

	return nil
}
