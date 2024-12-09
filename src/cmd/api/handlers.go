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
	return fmt.Errorf("nothing to show yet!")
}

func (s *Server) DevPath(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("nothing to show yet!")
}
