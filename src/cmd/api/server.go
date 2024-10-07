package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/anfego22/cyborgnerd/services/errors"
)

type Server struct {
	Tmpl *template.Template
}

func (s *Server) Start(port string) {
	http.HandleFunc("/", errors.ToStandar(s.Hello))
	fmt.Printf("Serving on http://localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
