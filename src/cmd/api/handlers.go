package api

import (
	"fmt"
	"net/http"
)

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) error {
	return s.Tmpl.ExecuteTemplate(w, "hello-world", nil)
}

func (s *Server) Construction(w http.ResponseWriter, r *http.Request) error {
	return s.Tmpl.ExecuteTemplate(w, "construction", nil)
}

func (s *Server) ArtistPath(w http.ResponseWriter, r *http.Request) error {
	var pages []string
	i := 1
	for i < 9 {
		pages = append(pages, fmt.Sprintf("/static/art/page_%d.webp", i))
		i++
	}
	return s.Tmpl.ExecuteTemplate(w, "artist-intro", pages)
}

func (s *Server) DevPath(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("nothing to show yet!")
}
