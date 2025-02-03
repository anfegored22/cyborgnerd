package api

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Tmpl *template.Template
	StaticFS fs.FS
}

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (s *Server) ToStandar(h HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {
			slog.Error(e.Error())
			s.Tmpl.ExecuteTemplate(w, "error", nil)
		}
	}
}

func (s *Server) Start(port string) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      s.routes(port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited cleanly")
}



func (s *Server) routes(port string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ToStandar(s.Hello))
	mux.HandleFunc("/artist-path", s.ToStandar(s.ArtistPath))
	mux.HandleFunc("/udemy-notes", s.ToStandar(s.UdemyNotes))
	mux.HandleFunc("/sketches", s.ToStandar(s.Sketches))
	mux.HandleFunc("/assets", s.ToStandar(s.Assets))
	mux.HandleFunc("/dev-path", s.ToStandar(s.DevPath))
	
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.StaticFS))))
	fmt.Printf("Serving on http://localhost:%s\n", port)

	return mux
}
