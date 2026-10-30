package api

import (
	"context"
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
	Tmpl     *template.Template
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

func (s *Server) Start(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverError := make(chan error, 1)
	go func() {
		slog.Info("Starting server", "address", addr)
		serverError <- srv.ListenAndServe()
	}()
	select {
	case err := <-serverError:
		if err != http.ErrServerClosed {
			return err
		}
		return nil
	case <-quit.Done():
	}

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Server exited cleanly")
	return nil
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ToStandar(s.Hello))
	mux.HandleFunc("/artist-path", s.ToStandar(s.ArtistPath))
	mux.HandleFunc("/udemy-notes", s.ToStandar(s.UdemyNotes))
	mux.HandleFunc("/sketches", s.ToStandar(s.Sketches))
	mux.HandleFunc("/assets", s.ToStandar(s.Assets))
	mux.HandleFunc("/dev-path", s.ToStandar(s.DevPath))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.StaticFS))))

	return mux
}
