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

func (s *Server) UdemyNotes(w http.ResponseWriter, r *http.Request) error {
	var pages []string
	i := 1
	for i < 9 {
		pages = append(pages, fmt.Sprintf("/static/art/page_%d.webp", i))
		i++
	}
	return s.Tmpl.ExecuteTemplate(w, "udemy-notes", pages)
}

func (s *Server) ArtistPath(w http.ResponseWriter, r *http.Request) error {
	cc := []InfoCard{
		{Title: "Udemy Notes", Description: "Detailed notes from the Udemy course, capturing every essential concept and technique.", URL: "/udemy-notes"},
		{Title: "Sketches", Description: "A collection of early video game concept art and rough drafts.", URL: "/sketches"},
		{Title: "Schedule", Description: "An overview of the planned timeline and milestones for the project.", URL: "/schedule"},
	}
	return s.Tmpl.ExecuteTemplate(w, "artist-path", cc)
}


func (s *Server) DevPath(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("nothing to show yet!")
}
