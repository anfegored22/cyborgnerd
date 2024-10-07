package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func hello(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "hello-world", nil)
}

func main() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}
