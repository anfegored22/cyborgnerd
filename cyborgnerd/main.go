package main

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed templates/*.html static/*
var content embed.FS

type page struct {
	Title string
	Path  string
}

func routes() http.Handler {
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))
	static, err := fs.Sub(content, "static")
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	for path, title := range map[string]string{"/": "Home", "/about": "About", "/projects": "Projects"} {
		view := page{Title: title, Path: path}
		pattern := "GET " + path
		if path == "/" {
			pattern = "GET /{$}"
		}
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Content-Language", "en")
			if err := tmpl.ExecuteTemplate(w, "page.html", view); err != nil {
				slog.Error("Render page", "error", err)
			}
		})
	}
	return mux
}

func run() error {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	server := &http.Server{
		Addr: addr, Handler: routes(),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errors := make(chan error, 1)
	go func() {
		slog.Info("Starting CyborgNerd", "address", addr)
		errors <- server.ListenAndServe()
	}()
	select {
	case err := <-errors:
		if err != http.ErrServerClosed {
			return err
		}
		return nil
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}

func main() {
	if err := run(); err != nil {
		slog.Error("Server stopped", "error", err)
		os.Exit(1)
	}
}
