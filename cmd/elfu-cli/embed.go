package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/wishmatic/elfu-mcp/internal/cli"
	"github.com/wishmatic/elfu-mcp/internal/present"
)

func newEmbedCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "embed <video-url>",
		Short: "Print the HTML player for a video URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return printEmbed(deps, args[0])
		},
	}
}

func printEmbed(deps cli.Deps, input string) error {
	videoURL := strings.TrimSpace(input)
	if videoURL == "" {
		return errors.New("embed: a video URL is required")
	}

	html, err := embeddedHTML(present.EmbedVideo(videoURL))
	if err != nil {
		return err
	}

	fmt.Fprintln(deps.Out, html)

	return nil
}

func embeddedHTML(content []mcp.Content) (string, error) {
	for _, item := range content {
		resource, ok := item.(*mcp.EmbeddedResource)
		if !ok || resource.Resource == nil {
			continue
		}

		return resource.Resource.Text, nil
	}

	return "", errors.New("embed: the video resource is missing")
}
