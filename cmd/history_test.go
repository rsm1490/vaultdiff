package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestHistoryCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "history <secret-path>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("history command not registered on root")
	}
}

func TestHistoryCmd_RequiresArg(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.AddCommand(historyCmd)

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"history"})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no secret-path provided")
	}
}

func TestHistoryCmd_MountFlag_Default(t *testing.T) {
	f := historyCmd.Flags().Lookup("mount")
	if f == nil {
		t.Fatal("mount flag not defined")
	}
	if f.DefValue != "secret" {
		t.Errorf("expected default mount 'secret', got '%s'", f.DefValue)
	}
}

func TestHistoryCmd_ShortDescription(t *testing.T) {
	if !strings.Contains(historyCmd.Short, "version") {
		t.Errorf("short description should mention 'version', got: %s", historyCmd.Short)
	}
}
