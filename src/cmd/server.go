package cmd

import (
	"html/template"
	"net/http"
)

type Server struct {
	Tmpl *template.Template
}

func (s *Server) Start() {
	http.HandleFunc("/", s.Hello)
	http.ListenAndServe(":8080", nil)
}
