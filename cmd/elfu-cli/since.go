package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/elfu-mcp/internal/cli"
	"github.com/wishmatic/elfu-mcp/internal/timing"
)

func newSinceCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "since <timestamp>",
		Short: "Report how long ago an ISO 8601 timestamp was",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return printSince(deps, args[0])
		},
	}
}

func printSince(deps cli.Deps, input string) error {
	then, err := timing.Parse(input, deps.Location)
	if err != nil {
		return fmt.Errorf("since: %w", err)
	}

	now := deps.Now().In(deps.Location)
	then = then.In(deps.Location)

	fmt.Fprintf(deps.Out, "Human: %s\n", timing.HumanSince(then, now))
	fmt.Fprintf(deps.Out, "At: %s\n", timing.Stamp{At: then}.ISO())
	fmt.Fprintf(deps.Out, "Seconds: %d\n", int64(now.Sub(then).Seconds()))
	fmt.Fprintf(deps.Out, "Now: %s\n", timing.Stamp{At: now}.ISO())

	return nil
}
