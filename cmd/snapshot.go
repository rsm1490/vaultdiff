package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var snapshotMount string
var snapshotOutput string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot <path> <version>",
		Short: "Capture a snapshot of a secret version",
		Args:  cobra.ExactArgs(2),
		RunE:  runSnapshot,
	}
	snapshotCmd.Flags().StringVar(&snapshotMount, "mount", "secret", "KV mount path")
	snapshotCmd.Flags().StringVar(&snapshotOutput, "output", "text", "Output format: text or json")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	sn := vault.NewSnapshotter(client)
	entry, err := sn.Capture(snapshotMount, path, version)
	if err != nil {
		return err
	}

	snap := vault.NewSnapshot(
		fmt.Sprintf("%s-v%d", path, version),
		[]vault.SnapshotEntry{*entry},
	)

	if snapshotOutput == "json" {
		return json.NewEncoder(os.Stdout).Encode(snap)
	}

	fmt.Fprintf(os.Stdout, "Snapshot: %s\n", snap.ID)
	fmt.Fprintf(os.Stdout, "Taken at: %s\n", snap.TakenAt.Format("2006-01-02T15:04:05Z"))
	for _, e := range snap.Entries {
		fmt.Fprintf(os.Stdout, "  path=%s version=%d keys=%d\n", e.Path, e.Version, len(e.Data))
	}
	return nil
}
