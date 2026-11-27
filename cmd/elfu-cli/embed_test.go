package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/wishmatic/elfu-mcp/internal/cli"
)

func TestEmbedPrintsThePlayerWithoutFetching(t *testing.T) {
	var requests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		http.Error(w, "should not be fetched", http.StatusInternalServerError)
	}))
	defer server.Close()

	videoURL := server.URL + "/video.mp4"

	out, _, err := runCommand(t, cli.Deps{}, "embed", videoURL)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if n := requests.Load(); n != 0 {
		t.Errorf("server received %d requests, want none", n)
	}

	for _, want := range []string{"<!DOCTYPE html>", `<video src="` + videoURL + `"`} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout does not contain %q:\n%s", want, out)
		}
	}
}

func TestEmbedEscapesTheURL(t *testing.T) {
	out, _, err := runCommand(t, cli.Deps{}, "embed", "https://example.com/video.mp4?a=1&b=2")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(out, `src="https://example.com/video.mp4?a=1&amp;b=2"`) {
		t.Errorf("stdout does not escape the URL:\n%s", out)
	}
}

func TestEmbedRejectsAnEmptyURL(t *testing.T) {
	out, _, err := runCommand(t, cli.Deps{}, "embed", "   ")
	if err == nil || !strings.Contains(err.Error(), "video URL is required") {
		t.Errorf("Execute() error = %v, want an empty URL rejection", err)
	}

	if out != "" {
		t.Errorf("stdout = %q, want nothing", out)
	}
}
