package api

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type Server struct {
	Tmpl *template.Template
}

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (s *Server) ToStandar(h HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {
			slog.Error(e.Error())
			s.Tmpl.ExecuteTemplate(w, "construction", nil)
		}
	}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/", s.ToStandar(s.Hello))
	http.HandleFunc("/artist-path", s.ToStandar(s.ArtistPath))
	http.HandleFunc("/dev-path", s.ToStandar(s.DevPath))

	fmt.Printf("Serving on http://localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
