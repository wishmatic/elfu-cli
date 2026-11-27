package main

import (
	"strings"
	"testing"

	"github.com/wishmatic/elfu-mcp/internal/cli"
)

func TestHelpListsEveryCommand(t *testing.T) {
	out, _, err := runCommand(t, cli.Deps{}, "help")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	for _, command := range []string{"serve", "time", "since", "sleep", "embed", "inline"} {
		if !strings.Contains(out, command) {
			t.Errorf("help output does not mention %q:\n%s", command, out)
		}
	}
}

func TestServeHelpDoesNotStartTheServer(t *testing.T) {
	out, _, err := runCommand(t, cli.Deps{}, "serve", "--help")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(out, "elfu-cli serve") {
		t.Errorf("help output = %q, want the serve usage", out)
	}
}
