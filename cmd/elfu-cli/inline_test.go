package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/elfu-mcp/internal/cli"
)

func imageServer(t *testing.T, mediaType string, body []byte) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", mediaType)
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)

	return server
}

func storedFiles(t *testing.T, dir string) []string {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%s) error: %v", dir, err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names
}

func TestInlineWritesThePreparedImage(t *testing.T) {
	server := imageServer(t, "image/png", pngBytes(t))
	dir := filepath.Join(t.TempDir(), "images")

	out, _, err := runCommand(t, cli.Deps{}, "inline", server.URL+"/logo.png", "--out", dir)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	names := storedFiles(t, dir)
	if len(names) != 1 {
		t.Fatalf("directory has %d files, want one", len(names))
	}

	path := filepath.Join(dir, names[0])

	if filepath.Ext(path) != ".webp" {
		t.Errorf("extension = %s, want .webp", filepath.Ext(path))
	}

	if !strings.Contains(out, "Path: "+path) {
		t.Errorf("stdout does not name %s:\n%s", path, out)
	}

	if !strings.Contains(out, "Media type: image/webp") {
		t.Errorf("stdout does not report the media type:\n%s", out)
	}

	if !strings.Contains(out, "Bytes: ") {
		t.Errorf("stdout does not report the size:\n%s", out)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}

	if !bytes.HasPrefix(data, []byte("RIFF")) || len(data) < 12 || string(data[8:12]) != "WEBP" {
		t.Error("stored file is not a WebP")
	}
}

func TestInlineWritesToTheCurrentDirectoryByDefault(t *testing.T) {
	t.Chdir(t.TempDir())

	server := imageServer(t, "image/png", pngBytes(t))

	out, _, err := runCommand(t, cli.Deps{}, "inline", server.URL)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	names := storedFiles(t, ".")
	if len(names) != 1 {
		t.Fatalf("directory has %d files, want one", len(names))
	}

	if !strings.Contains(out, "Path: "+names[0]) {
		t.Errorf("stdout does not name %s:\n%s", names[0], out)
	}
}

func TestInlineSurfacesAFetchFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "gone", http.StatusInternalServerError)
	}))
	defer server.Close()

	dir := t.TempDir()

	_, _, err := runCommand(t, cli.Deps{}, "inline", server.URL, "--out", dir)
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Errorf("Execute() error = %v, want the upstream status", err)
	}

	if names := storedFiles(t, dir); len(names) != 0 {
		t.Errorf("directory has %v, want nothing written", names)
	}
}

func TestInlineSurfacesAnUndecodableImage(t *testing.T) {
	server := imageServer(t, "image/png", brokenPNG)
	dir := t.TempDir()

	_, _, err := runCommand(t, cli.Deps{}, "inline", server.URL, "--out", dir)
	if err == nil || !strings.Contains(err.Error(), "imgfmt") {
		t.Errorf("Execute() error = %v, want a decode failure", err)
	}

	if names := storedFiles(t, dir); len(names) != 0 {
		t.Errorf("directory has %v, want nothing written", names)
	}
}
