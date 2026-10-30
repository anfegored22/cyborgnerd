package main

import (
	"io/fs"
	"strings"
	"testing"
)

func TestEmbeddedAssets(t *testing.T) {
	for _, name := range []string{
		"templates/home.html",
		"static/portfolio.css",
		"static/js/portfolio.js",
		"static/fonts/Pixelion-Regular.woff2",
		"static/portfolio/pages/page-01.webp",
		"static/portfolio/web/page-10-image-052.webp",
	} {
		if _, err := fs.Stat(content, name); err != nil {
			t.Errorf("missing deployable asset %s: %v", name, err)
		}
	}
	err := fs.WalkDir(content, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".pdf") || strings.HasPrefix(path, "static/portfolio/page-") {
			t.Errorf("private source asset embedded: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
