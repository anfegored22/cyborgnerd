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
	pages := []string{"page_1.png", "page_2.png", "page_3.png", "page_4.png", "page_5.png"}
	return s.Tmpl.ExecuteTemplate(w, "artist-intro", pages)
}

func (s *Server) DevPath(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("nothing to show yet!")
}
