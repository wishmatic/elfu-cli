# `elfu-cli` Plan

## Goal

Add an `elfu-cli` binary that exposes the MCP tools' behavior without the HTTP server or an MCP
client, the way `nano-mcp` ships `cmd/nano-cli`. `serve` runs the MCP server itself, so `cmd/server`
is replaced rather than kept alongside.

The command surface is:

```text
elfu-cli help
elfu-cli serve
elfu-cli time
elfu-cli since <timestamp>
elfu-cli sleep <seconds>
elfu-cli embed <video-url>
elfu-cli inline <image-url>
```

## Scope

In scope:

- `cmd/elfu-cli` as the only entry point, replacing `cmd/server`.
- An `internal/cli` package for behavior shared across commands.
- One command per MCP tool, plus `serve`.
- Unit and command-level tests, README and Dockerfile updates.

Out of scope:

- Changing the MCP tool surface or server behavior.
- Credential storage, `login`, and `logout`: the tools need no upstream credentials, and `API_KEY`
  stays the server's own request authentication.
- Publishing prebuilt binaries, a TUI, or shell completion installation.

## Decisions

- `github.com/spf13/cobra`, as in `nano-mcp`, with `SilenceUsage` and `SilenceErrors` on the root so
  a runtime failure prints one line.
- `internal/timing` stays the single source for timestamp formatting and humanized durations, so
  the CLI and the tools report the same text and fields.
- The CLI reports the same values the tools put in their structured output, as readable lines, so a
  shell user sees what an agent would.
- `inline` writes the prepared WebP to a local file rather than returning an image block, because a
  CLI has no content blocks: `--out` selects the directory (default `.`), and the command prints the
  path, media type, and size. The image is read through `internal/resolve` with
  `resolve.New(nil, "", sources)`, so http(s) URLs, data URIs, and raw base64 work.
- The local timezone stays the process's own, which is what the server's `TZ` already sets.

## Design

### Package layout

```text
cmd/elfu-cli/main.go          root command and process entry point
cmd/elfu-cli/serve.go         MCP server command, logger, signal handling
cmd/elfu-cli/time.go          time command
cmd/elfu-cli/since.go         since command
cmd/elfu-cli/sleep.go         sleep command
cmd/elfu-cli/embed.go         embed command
cmd/elfu-cli/inline.go        inline command
internal/cli/deps.go          dependency injection, config, resolver, clock
internal/cli/asset.go         local image writing
```

`cmd/elfu-cli` may import `internal/{config,imgfmt,resolve,sourcemap,server,timing}`. It must not
import `internal/mcp`, which stays the presentation layer for MCP content.

### Dependency injection

`internal/cli.Deps` carries `In`, `Out`, `Err`, `Config config.Config`, `Resolver *resolve.Resolver`,
`Location *time.Location`, and `Now func() time.Time`. Nil fields take the process defaults
(`time.Local`, `time.Now`), so tests can pin the clock and capture output in a buffer.

`internal/cli.WriteAsset(dir, prefix string, data []byte, mediaType string) (string, error)` picks
the extension from the media type, creates the directory, writes through a temporary file, and
renames it into place.

### Commands

`serve` is the old `cmd/server/main.go` moved behind Cobra: load config, build the logger, run
`internal/server`, and shut down gracefully on `SIGINT`/`SIGTERM` with the same 15 second window.
The `_ "time/tzdata"` import moves with it.

`time` prints the same fields the tool returns as one line: the ISO timestamp, the zone and offset,
the Unix seconds, and the readable form.

`since <timestamp>` prints the humanized duration, the resolved timestamp, the seconds, and the now
it was measured against. A timestamp without a zone is read in the local timezone, as the tool does.

`sleep <seconds>` waits up to 30 seconds, refusing more, and prints the time it woke.

`embed <video-url>` prints the same HTML the tool returns, without downloading anything.

`inline <image-url>` fetches the image, runs `imgfmt.Inline` with the same edge and byte limits, and
writes the result, printing the path, media type, and byte count. A failure to prepare the image
exits non-zero with the reason.

### Errors and exit codes

Success exits `0`; usage and runtime failures exit `1` with one line on stderr. No partial file is
written: the asset is written only after the image was prepared.

## Implementation Units

### Unit 1: Scaffold, `serve`, and shared helpers

Add Cobra, `cmd/elfu-cli/main.go`, `cmd/elfu-cli/serve.go`, `internal/cli`, and move the entry point
out of `cmd/server`. Update the Dockerfile to build `cmd/elfu-cli`, copy it to `/elfu-cli`, and use
`["/elfu-cli", "serve"]` as the entry point.

Acceptance criteria:

- [x] `go build ./cmd/elfu-cli` produces a binary whose `elfu-cli help` exits 0 with the tree above.
- [x] `cmd/server` no longer exists and nothing in the repo references it.
- [x] The Dockerfile builds `./cmd/elfu-cli` and its entry point is `elfu-cli serve`.
- [x] `elfu-cli serve` loads `.env`/environment config, keeps the embedded zoneinfo, and shuts down
      gracefully on `SIGINT`.
- [x] `internal/cli` has tests for injected deps and asset writing (extension, atomic rename,
      unknown media type, missing directory creation).

### Unit 2: `time`, `since`, and `sleep`

Acceptance criteria:

- [x] `time` prints the ISO timestamp, zone, offset, Unix seconds, and readable form for a pinned
      clock and location.
- [x] `since` prints the humanized duration and seconds for past and future timestamps, reads a
      zone-less timestamp in the local zone, and rejects an unparseable one without output.
- [x] `sleep` waits the requested duration, refuses more than 30 seconds and negative values before
      waiting, and prints the wake time.
- [x] Tests cover all three with an injected clock.

### Unit 3: `embed` and `inline`

Acceptance criteria:

- [x] `embed` prints the HTML player for the URL and does not make a network request.
- [x] `embed` rejects an empty URL.
- [x] `inline` fetches through the resolver, prepares the image with `imgfmt.Inline`, and writes it
      to `--out` with the extension of its media type, printing path, media type, and size.
- [x] `inline` surfaces a fetch failure and an undecodable image without writing a file.
- [x] Tests cover both commands against `httptest` servers.

### Unit 4: Documentation

Acceptance criteria:

- [x] `README.md` documents `go install github.com/wishmatic/elfu-mcp/cmd/elfu-cli@latest`, the
      commands above, and that `inline` writes locally while `serve` needs `API_KEY` and
      `PUBLIC_HOST`.
- [x] The Docker and local-development examples use `elfu-cli serve`.
- [x] `go build ./...`, `go vet ./...`, and `go test ./... -race -count=1` pass.

## Verification Commands

```sh
go build ./...
go vet ./...
go test ./... -race -count=1
go build -o /tmp/elfu-cli ./cmd/elfu-cli
/tmp/elfu-cli help
/tmp/elfu-cli time
/tmp/elfu-cli since 2026-01-01
/tmp/elfu-cli serve --help
```

Human check: `elfu-cli inline https://example.com/logo.webp --out ./images` against a real URL, then
open the written file.
