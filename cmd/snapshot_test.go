package cmd

import (
	"testing"
)

func TestSnapshotCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "snapshot <path> <version>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected snapshot command to be registered on root")
	}
}

func TestSnapshotCmd_RequiresTwoArgs(t *testing.T) {
	cmd := findCommand(rootCmd, "snapshot")
	if cmd == nil {
		t.Fatal("snapshot command not found")
	}
	if err := cmd.Args(cmd, []string{"only-one"}); err == nil {
		t.Error("expected error when only one arg provided")
	}
}

func TestSnapshotCmd_MountFlag_Default(t *testing.T) {
	cmd := findCommand(rootCmd, "snapshot")
	if cmd == nil {
		t.Fatal("snapshot command not found")
	}
	f := cmd.Flags().Lookup("mount")
	if f == nil {
		t.Fatal("expected mount flag")
	}
	if f.DefValue != "secret" {
		t.Errorf("expected default mount 'secret', got %s", f.DefValue)
	}
}

func TestSnapshotCmd_OutputFlag_Default(t *testing.T) {
	cmd := findCommand(rootCmd, "snapshot")
	if cmd == nil {
		t.Fatal("snapshot command not found")
	}
	f := cmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("expected output flag")
	}
	if f.DefValue != "text" {
		t.Errorf("expected default output 'text', got %s", f.DefValue)
	}
}

func TestSnapshotCmd_ShortDescription(t *testing.T) {
	cmd := findCommand(rootCmd, "snapshot")
	if cmd == nil {
		t.Fatal("snapshot command not found")
	}
	if cmd.Short == "" {
		t.Error("expected non-empty short description")
	}
}
