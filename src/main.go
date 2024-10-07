package main

import (
	"embed"
	"html/template"

	"github.com/anfego22/cyborgnerd/cmd/api"
)

//go:embed templates/*
var content embed.FS

func main() {
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))
	cn := api.Server{Tmpl: tmpl}
	cn.Start("8080")

}
