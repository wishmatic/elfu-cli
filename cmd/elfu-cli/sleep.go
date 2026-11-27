package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/wishmatic/elfu-mcp/internal/cli"
	"github.com/wishmatic/elfu-mcp/internal/timing"
)

const maxSleepSeconds = 30

func newSleepCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "sleep <seconds>",
		Short: "Wait for a number of seconds, up to 30, then report the wake time",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSleep(cmd.Context(), deps, args[0])
		},
	}
}

func runSleep(ctx context.Context, deps cli.Deps, raw string) error {
	seconds, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return fmt.Errorf("sleep: %q is not a number of seconds", raw)
	}

	wait, err := sleepDuration(seconds)
	if err != nil {
		return err
	}

	if err := waitFor(ctx, wait); err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	fmt.Fprintf(deps.Out, "Seconds: %g\n", seconds)
	fmt.Fprintf(deps.Out, "Now: %s\n", timing.Stamp{At: deps.Now().In(deps.Location)}.ISO())

	return nil
}

func sleepDuration(seconds float64) (time.Duration, error) {
	if math.IsNaN(seconds) || seconds < 0 {
		return 0, fmt.Errorf("sleep: seconds must be zero or more, got %v", seconds)
	}

	if seconds > maxSleepSeconds {
		return 0, fmt.Errorf("sleep: seconds must be at most %d, got %v", maxSleepSeconds, seconds)
	}

	return time.Duration(seconds * float64(time.Second)), nil
}

func waitFor(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
