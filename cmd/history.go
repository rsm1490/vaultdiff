package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	historyMount string
)

var historyCmd = &cobra.Command{
	Use:   "history <secret-path>",
	Short: "List all versions of a secret with metadata",
	Args:  cobra.ExactArgs(1),
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().StringVar(&historyMount, "mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	secretPath := args[0]

	client, err := newVaultClient()
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	ctx := context.Background()
	history, err := client.GetHistory(ctx, historyMount, secretPath)
	if err != nil {
		return fmt.Errorf("fetching history: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tCREATED\tDESTROYED\tDELETED")

	for _, v := range history.Versions {
		destroyed := "-"
		if v.Destroyed {
			destroyed = "yes"
		}
		deleted := "-"
		if !v.DeletedAt.IsZero() {
			deleted = v.DeletedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			v.Version,
			v.CreatedAt.Format("2006-01-02 15:04:05"),
			destroyed,
			deleted,
		)
	}

	return w.Flush()
}
