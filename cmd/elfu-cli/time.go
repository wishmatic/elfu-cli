package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/elfu-mcp/internal/cli"
	"github.com/wishmatic/elfu-mcp/internal/timing"
)

func newTimeCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "time",
		Short: "Report the current date and time",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			return printTime(deps)
		},
	}
}

func printTime(deps cli.Deps) error {
	stamp := timing.Stamp{At: deps.Now().In(deps.Location)}

	fmt.Fprintf(deps.Out, "%s (%s %s), Unix %d, %s\n",
		stamp.ISO(), stamp.Zone(), stamp.Offset(), stamp.Unix(), stamp.Readable(),
	)

	return nil
}
