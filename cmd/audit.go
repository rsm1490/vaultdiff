package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/audit"
	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

var auditLogFile string

var auditCmd = &cobra.Command{
	Use:   "audit <path> <version-a> <version-b>",
	Short: "Diff two secret versions and append the result to an audit log",
	Args:  cobra.ExactArgs(3),
	RunE:  runAudit,
}

func init() {
	auditCmd.Flags().StringVarP(&auditLogFile, "log-file", "l", "vaultdiff-audit.log", "path to the audit log file")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	path := args[0]
	versionA, err := parseVersion(args[1], "version-a")
	if err != nil {
		return err
	}
	versionB, err := parseVersion(args[2], "version-b")
	if err != nil {
		return err
	}

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secA, err := client.GetSecretVersion(path, versionA)
	if err != nil {
		return fmt.Errorf("fetch version %d: %w", versionA, err)
	}
	secB, err := client.GetSecretVersion(path, versionB)
	if err != nil {
		return fmt.Errorf("fetch version %d: %w", versionB, err)
	}

	changes := diff.Compare(secA.Data, secB.Data)
	summary := diff.Summarize(changes)

	logger, f, err := audit.NewFileLogger(auditLogFile)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := audit.Entry{
		Path:     path,
		VersionA: versionA,
		VersionB: versionB,
		Changes:  changes,
		Summary:  summary,
	}
	if err := logger.Record(entry); err != nil {
		return fmt.Errorf("write audit log: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Audit entry written to %s\n", auditLogFile)
	fmt.Fprintf(os.Stdout, "Changes: +%d -%d ~%d\n", summary.Added, summary.Removed, summary.Modified)
	return nil
}

// parseVersion converts a string argument into a positive integer version number.
func parseVersion(s, name string) (int, error) {
	var v int
	if _, err := fmt.Sscan(s, &v); err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", name, s, err)
	}
	if v <= 0 {
		return 0, fmt.Errorf("invalid %s %q: must be a positive integer", name, s)
	}
	return v, nil
}
