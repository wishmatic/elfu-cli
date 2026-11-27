package cli

import (
	"os"
	"testing"
	"time"

	"github.com/wishmatic/elfu-mcp/internal/config"
	"github.com/wishmatic/elfu-mcp/internal/resolve"
)

func TestWithDefaultsFillsTheProcessDefaults(t *testing.T) {
	deps := WithDefaults(Deps{})

	if deps.In != os.Stdin {
		t.Error("In is not os.Stdin")
	}

	if deps.Out != os.Stdout {
		t.Error("Out is not os.Stdout")
	}

	if deps.Err != os.Stderr {
		t.Error("Err is not os.Stderr")
	}

	if deps.Location != time.Local {
		t.Errorf("Location = %v, want time.Local", deps.Location)
	}

	if deps.Now == nil {
		t.Error("Now is nil")
	}
}

func TestWithDefaultsKeepsInjectedValues(t *testing.T) {
	loc := time.FixedZone("AEST", 10*60*60)
	pinned := time.Date(2026, 11, 27, 10, 0, 0, 0, loc)

	got := WithDefaults(Deps{
		Location: loc,
		Now:      func() time.Time { return pinned },
	})

	if got.Location != loc {
		t.Errorf("Location = %v, want the injected zone", got.Location)
	}

	if !got.Now().Equal(pinned) {
		t.Errorf("Now() = %v, want %v", got.Now(), pinned)
	}
}

func TestImageResolverUsesTheInjectedResolver(t *testing.T) {
	injected, err := resolve.New(nil, "", nil)
	if err != nil {
		t.Fatalf("resolve.New() error: %v", err)
	}

	got, err := ImageResolver(Deps{Resolver: injected})
	if err != nil {
		t.Fatalf("ImageResolver() error: %v", err)
	}

	if got != injected {
		t.Error("ImageResolver() did not return the injected resolver")
	}
}

func TestImageResolverBuildsFromTheConfig(t *testing.T) {
	deps := Deps{Config: config.Config{
		ImageURLMap: "https://chat.example.com/images/=/data/librechat-data",
	}}

	got, err := ImageResolver(deps)
	if err != nil {
		t.Fatalf("ImageResolver() error: %v", err)
	}

	if got == nil {
		t.Error("ImageResolver() = nil, want a resolver")
	}
}

func TestImageResolverReportsAnInvalidMap(t *testing.T) {
	_, err := ImageResolver(Deps{Config: config.Config{ImageURLMap: "not-a-pair"}})
	if err == nil {
		t.Fatal("ImageResolver() error = nil, want a map rejection")
	}
}

func TestDefaultDepsReadsTheEnvironment(t *testing.T) {
	t.Setenv("PORT", "9090")

	deps, err := DefaultDeps()
	if err != nil {
		t.Fatalf("DefaultDeps() error: %v", err)
	}

	if deps.Config.Port != 9090 {
		t.Errorf("Port = %d, want 9090", deps.Config.Port)
	}

	if deps.Out == nil {
		t.Error("Out is nil")
	}
}

func TestDefaultDepsReportsAnInvalidEnvironment(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	if _, err := DefaultDeps(); err == nil {
		t.Error("DefaultDeps() error = nil, want a parse failure")
	}
}
