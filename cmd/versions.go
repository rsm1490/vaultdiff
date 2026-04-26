package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultdiff/internal/vault"
)

var versionsMount string

var versionsCmd = &cobra.Command{
	Use:   "versions <secret-path>",
	Short: "List all versions of a Vault KVv2 secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runVersions,
}

func init() {
	versionsCmd.Flags().StringVar(&versionsMount, "mount", "secret", "KVv2 mount path")
	rootCmd.AddCommand(versionsCmd)
}

func runVersions(cmd *cobra.Command, args []string) error {
	secretPath := args[0]

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	versions, err := client.ListVersions(context.Background(), versionsMount, secretPath)
	if err != nil {
		return fmt.Errorf("listing versions: %w", err)
	}

	if len(versions) == 0 {
		fmt.Println("No versions found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tCREATED\tDELETION TIME\tDESTROYED")
	for _, v := range versions {
		destroyed := "-"
		if v.Destroyed {
			destroyed = "yes"
		}
		deletionTime := v.DeletionTime
		if deletionTime == "" {
			deletionTime = "-"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			v.Version, v.CreatedTime, deletionTime, destroyed)
	}
	return w.Flush()
}
