package main

import (
	"embed"
	"html/template"

	"github.com/anfego22/cyborgnerd/cmd"
)

//go:embed templates/*
var content embed.FS

func main() {
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))
	cn := cmd.Server{Tmpl: tmpl}
	cn.Start()
}
