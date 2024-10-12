package main

import (
	"html/template"
	"net/http"

	"github.com/anfego22/cyborgnerd/cmd/api"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	cn := api.Server{Tmpl: tmpl}
	cn.Start("8080")
}
