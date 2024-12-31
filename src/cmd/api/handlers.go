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
	pages := ImagesPath("/static/art", 1, 9)
	data := struct {
		Title string
		Pages []string
	}{Title: "Udemy Notes", Pages: pages}
	return s.Tmpl.ExecuteTemplate(w, "card-expanded", data)
}

func (s *Server) Sketches(w http.ResponseWriter, r *http.Request) error {
	pages := ImagesPath("/static/tima", 0, 10)
	data := struct {
		Title string
		Pages []string
	}{Title: "Sketches", Pages: pages}
	return s.Tmpl.ExecuteTemplate(w, "card-expanded", data)
}

func ImagesPath(path string, i0, n int) []string{
	var pages []string
	i := i0
	for i < n {
		pages = append(pages, fmt.Sprintf("%s/page_%d.webp", path, i))
		i++
	}
	return pages
}

func (s *Server) ArtistPath(w http.ResponseWriter, r *http.Request) error {
	cc := []InfoCard{
		{Title: "Udemy Notes", Description: "Detailed notes from the Udemy course, capturing every essential concept and technique.", URL: "/udemy-notes"},
		{Title: "Sketches", Description: "A collection of early video game concept art and rough drafts.", URL: "/sketches"},
	}
	if r.Header.Get("HX-Request") == "true" {
		return s.Tmpl.ExecuteTemplate(w, "artist-path", cc) 
	}
	return s.Tmpl.ExecuteTemplate(w, "full-artist-path", cc)
}


func (s *Server) DevPath(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("nothing to show yet!")
}
