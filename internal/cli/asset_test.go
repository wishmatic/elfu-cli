package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteAssetWritesTheFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "output")

	path, err := WriteAsset(dir, "elfu-image", []byte("image bytes"), "image/webp")
	if err != nil {
		t.Fatalf("WriteAsset() error: %v", err)
	}

	if filepath.Dir(path) != dir {
		t.Errorf("path = %s, want it under %s", path, dir)
	}

	if filepath.Ext(path) != ".webp" {
		t.Errorf("extension = %s, want .webp", filepath.Ext(path))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}

	if string(data) != "image bytes" {
		t.Errorf("data = %q, want image bytes", data)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0o644 {
		t.Errorf("mode = %o, want 644", perm)
	}
}

func TestWriteAssetLeavesNoTemporaryFile(t *testing.T) {
	dir := t.TempDir()

	if _, err := WriteAsset(dir, "elfu-image", []byte("image bytes"), "image/png"); err != nil {
		t.Fatalf("WriteAsset() error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("directory has %d entries, want only the stored file", len(entries))
	}

	if name := entries[0].Name(); !strings.HasSuffix(name, ".png") {
		t.Errorf("stored file = %s, want a .png", name)
	}
}

func TestWriteAssetUsesTheMediaTypeExtension(t *testing.T) {
	for mediaType, extension := range map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/jxl":  ".jxl",
		"image/webp": ".webp",
	} {
		t.Run(mediaType, func(t *testing.T) {
			path, err := WriteAsset(t.TempDir(), "elfu-image", []byte("data"), mediaType)
			if err != nil {
				t.Fatalf("WriteAsset() error: %v", err)
			}

			if filepath.Ext(path) != extension {
				t.Errorf("extension = %s, want %s", filepath.Ext(path), extension)
			}
		})
	}
}

func TestWriteAssetRejectsAnUnknownMediaType(t *testing.T) {
	_, err := WriteAsset(t.TempDir(), "elfu-image", []byte("data"), "application/pdf")
	if err == nil || !strings.Contains(err.Error(), "cannot store") {
		t.Errorf("WriteAsset() error = %v, want a media-type rejection", err)
	}
}

func TestWriteAssetNamesTheFileWithItsPrefix(t *testing.T) {
	path, err := WriteAsset(t.TempDir(), "elfu-image", []byte("data"), "image/png")
	if err != nil {
		t.Fatalf("WriteAsset() error: %v", err)
	}

	if !strings.HasPrefix(filepath.Base(path), "elfu-image-") {
		t.Errorf("name = %s, want the prefix", filepath.Base(path))
	}
}
