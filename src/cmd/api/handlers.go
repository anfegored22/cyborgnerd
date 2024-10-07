package api

import "net/http"

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) error {
	return s.Tmpl.ExecuteTemplate(w, "hello-world", nil)
}
