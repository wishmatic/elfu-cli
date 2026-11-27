package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wishmatic/elfu-mcp/internal/cli"
)

func main() {
	deps, err := cli.DefaultDeps()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := newRoot(deps).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRoot(deps cli.Deps) *cobra.Command {
	deps = cli.WithDefaults(deps)

	root := &cobra.Command{
		Use:           "elfu-cli",
		Short:         "Use Elfu's small utilities from the command line",
		Long:          "Report the time, wait, embed videos, and inline images, without an MCP client.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.SetIn(deps.In)
	root.SetOut(deps.Out)
	root.SetErr(deps.Err)

	root.AddCommand(
		newServeCommand(),
		newTimeCommand(deps),
		newSinceCommand(deps),
		newSleepCommand(deps),
		newEmbedCommand(deps),
		newInlineCommand(deps),
	)

	return root
}
