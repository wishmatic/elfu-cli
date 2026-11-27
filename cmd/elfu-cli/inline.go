package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/elfu-mcp/internal/cli"
	"github.com/wishmatic/elfu-mcp/internal/imgfmt"
)

func newInlineCommand(deps cli.Deps) *cobra.Command {
	var out string

	cmd := &cobra.Command{
		Use:   "inline <image-url>",
		Short: "Fetch an image, prepare it as WebP, and write it locally",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInline(cmd.Context(), deps, args[0], out)
		},
	}

	cmd.Flags().StringVar(&out, "out", ".", "directory for the prepared image")

	return cmd
}

func runInline(ctx context.Context, deps cli.Deps, imageURL, out string) error {
	resolver, err := cli.ImageResolver(deps)
	if err != nil {
		return err
	}

	source, err := resolver.Resolve(ctx, imageURL)
	if err != nil {
		return fmt.Errorf("inline: %w", err)
	}

	prepared, err := imgfmt.Inline(source.Data, imgfmt.InlineMaxEdge, imgfmt.InlineMaxBytes)
	if err != nil {
		return fmt.Errorf("inline: %w", err)
	}

	path, err := cli.WriteAsset(out, "elfu-image", prepared.Data, prepared.MediaType)
	if err != nil {
		return fmt.Errorf("inline: %w", err)
	}

	fmt.Fprintf(deps.Out, "Path: %s\n", path)
	fmt.Fprintf(deps.Out, "Media type: %s\n", prepared.MediaType)
	fmt.Fprintf(deps.Out, "Bytes: %d\n", len(prepared.Data))

	return nil
}
