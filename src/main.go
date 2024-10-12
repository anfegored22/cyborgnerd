package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/anfego22/cyborgnerd/cmd/api"
)

//go:embed templates/* static/*
var content embed.FS

func main() {
	staticFS, _ := fs.Sub(content, "static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))
	cn := api.Server{Tmpl: tmpl}
	cn.Start("8080")

}
