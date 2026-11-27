package main

import (
	"bytes"
	"image"
	"image/png"
	"testing"
	"time"

	"github.com/wishmatic/elfu-mcp/internal/cli"
)

var testLocation = time.FixedZone("AEST", 10*60*60)

var testNow = time.Date(2026, 11, 27, 10, 0, 0, 0, testLocation)

var brokenPNG = append([]byte("\x89PNG\r\n\x1a\n"), []byte("garbage")...)

func runCommand(t *testing.T, deps cli.Deps, args ...string) (string, string, error) {
	t.Helper()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	deps.Out = out
	deps.Err = errOut

	command := newRoot(deps)
	command.SetArgs(args)

	err := command.Execute()

	return out.String(), errOut.String(), err
}

func pinnedDeps() cli.Deps {
	return cli.Deps{
		Location: testLocation,
		Now:      func() time.Time { return testNow },
	}
}

func pngBytes(t *testing.T) []byte {
	t.Helper()

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, image.NewRGBA(image.Rect(0, 0, 8, 8))); err != nil {
		t.Fatalf("png.Encode() error: %v", err)
	}

	return buf.Bytes()
}
