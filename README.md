# Elfu CLI and MCP

> I no longer use Librechat, and my new agent of choice can do these things without help, so this
> repository has been archived.

<img src="docs/images/logo.webp" alt="Elfu MCP Logo" width="128">

Formerly "Librechat Functional Utilities," LFU, or Elfu. Now, it's general utilities that can help
coding and general assistant AI agents.

## Tools

### `inline`

`inline` takes an `image_url`, follows redirects, and returns the image, downscaled and re-encoded
to WebP.

This is useful for cases where the model has vision capability but only a URL. After inlining, in
the same turn, a Librechat agent/LLM can see the image.

#### What's the point of the `IMAGE_URL_MAP`?

If you are running other MCP that serve files but themselves are internally networked or not
accessible by this MCP because it's on a reverse proxy or some other reason, you need a way to
translate those "public" URLs to internal ones. This MCP defines a mapping environment variable that
does this.

This also works with local images on disk if mapped, though you should be careful here; the MCP can
enumerate and see all files uploaded by any user to Librechat, for example, so make sure URL paths
are unguessable.

### `embed`

`embed` takes a `video_url` and returns an HTML player that Librechat renders inline in the message,
so the user can watch the video in the chat.

It downloads nothing; the user's browser loads the URL itself. That URL therefore has to be one the
browser can load on its own, without the chat's credentials. A video only reachable from inside the
network, such as one behind Librechat's authenticated `/images/*` route, will not play.

### `time`, `since`, and `sleep`

`time` returns the datetime. `since` takes a datetime and returns how long ago it was in
human-readable form. `sleep` waits up to 30 seconds and returns the time it woke.

## Usage

### CLI

```sh
go install github.com/wishmatic/elfu-mcp/cmd/elfu-cli@latest
```

Examples:

```sh
elfu-cli help
elfu-cli time
elfu-cli since 2026-01-01
elfu-cli sleep 2
elfu-cli embed https://example.com/video.mp4
elfu-cli inline https://example.com/logo.webp --out ./images
```

`embed` prints the same HTML player the tool returns. `inline` writes the prepared WebP to a local
file rather than returning an image block, since a CLI has no content blocks; `--out` picks the
directory.

### MCP server

The Docker image's entry point is `elfu-cli serve`, which needs `API_KEY` and `PUBLIC_HOST`:

```sh
docker run -d \
  -p 8080:8080 \
  -e API_KEY=change-me \
  -e PUBLIC_HOST=http://192.168.1.10:8080 \
  -e IMAGE_URL_MAP=https://chat.example.com/images/=/data/librechat-data \
  -e TZ=Australia/Sydney \
  -v /mnt/user/appdata/elfu-mcp:/data \
  -v /mnt/user/appdata/librechat/client/public/images:/data/librechat-data:ro \
  ghcr.io/wishmatic/elfu-mcp:latest
```

The MCP endpoint is served at `/mcp`; stored images are served from `/i/`.

See [.env.example](.env.example) for all configuration.

### Authentication

`API_KEY` is required on every `/mcp` request, sent as `Authorization: Bearer <API_KEY>`. Stored
images are served without authentication.

## Local Development

Copy `.env.example` to `.env` and fill it in, then:

```sh
go run ./cmd/elfu-cli serve
```

## License

Elfu MCP is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
