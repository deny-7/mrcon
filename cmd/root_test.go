package cmd

import (
	"bytes"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newRootCmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--version"})
	SetVersion("1.2.3", "abc123", "2025-06-29T12:00:00Z")
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if want := "mrcon version: 1.2.3"; !bytes.Contains([]byte(out), []byte(want)) {
		t.Errorf("expected version output, got: %s", out)
	}
}

func TestMissingRequiredFlags(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newRootCmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing required flags, got nil")
	}
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("Error: --host, --port, --password, and a command are required.")) {
		t.Errorf("expected error message, got: %s", out)
	}
}

// More tests can be added for silent, raw, and terminal modes with mocks.
