package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"os"

	"github.com/anfego22/cyborgnerd/tacyborg/cmd/api"
)

// Only embed deployable assets, never source artwork or purchase archives.
//
//go:embed templates/*.html static/*.css static/js static/fonts static/img static/art static/tima static/portfolio/web static/portfolio/pages static/portfolio/*.json static/portfolio/*.txt
var content embed.FS

func main() {
	staticFS, err := fs.Sub(content, "static")
	if err != nil {
		slog.Error("Load static files", "error", err)
		os.Exit(1)
	}
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))
	cn := api.Server{Tmpl: tmpl, StaticFS: staticFS}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	if err := cn.Start(addr); err != nil {
		slog.Error("Server stopped", "error", err)
		os.Exit(1)
	}
}
