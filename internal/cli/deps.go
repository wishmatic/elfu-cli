package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/wishmatic/elfu-mcp/internal/config"
	"github.com/wishmatic/elfu-mcp/internal/resolve"
	"github.com/wishmatic/elfu-mcp/internal/sourcemap"
)

type Deps struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer

	Config config.Config

	Resolver *resolve.Resolver
	Location *time.Location
	Now      func() time.Time
}

func DefaultDeps() (Deps, error) {
	cfg, err := config.Load()
	if err != nil {
		return Deps{}, err
	}

	return WithDefaults(Deps{Config: cfg}), nil
}

func WithDefaults(deps Deps) Deps {
	if deps.In == nil {
		deps.In = os.Stdin
	}

	if deps.Out == nil {
		deps.Out = os.Stdout
	}

	if deps.Err == nil {
		deps.Err = os.Stderr
	}

	if deps.Location == nil {
		deps.Location = time.Local
	}

	if deps.Now == nil {
		deps.Now = time.Now
	}

	return deps
}

func ImageResolver(deps Deps) (*resolve.Resolver, error) {
	if deps.Resolver != nil {
		return deps.Resolver, nil
	}

	sources, err := sourcemap.Parse(deps.Config.ImageURLMap)
	if err != nil {
		return nil, fmt.Errorf("IMAGE_URL_MAP: %w", err)
	}

	return resolve.New(nil, "", sources)
}
