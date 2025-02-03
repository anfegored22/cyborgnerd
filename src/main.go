package main

import (
	"embed"
	"html/template"
	"io/fs"

	"github.com/anfego22/cyborgnerd/cmd/api"
)

//go:embed templates/* static/*
var content embed.FS

func main() {
	staticFS, _ := fs.Sub(content, "static")
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))
	cn := api.Server{Tmpl: tmpl, StaticFS: staticFS}
	cn.Start("8080")
}
