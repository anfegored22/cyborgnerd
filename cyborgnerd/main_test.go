package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPages(t *testing.T) {
	handler := routes()
	for _, path := range []string{"/", "/about", "/projects"} {
		for _, boosted := range []bool{false, true} {
			r := httptest.NewRequest(http.MethodGet, path, nil)
			if boosted {
				r.Header.Set("HX-Request", "true")
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			body := w.Body.String()
			if w.Code != http.StatusOK {
				t.Fatalf("%s: status %d", path, w.Code)
			}
			for _, want := range []string{"<!doctype html>", "CyborgNerd", `id="content"`, `hx-boost="true"`, `hx-select="#content"`, `href="` + path + `" aria-current="page"`} {
				if !strings.Contains(body, want) {
					t.Errorf("%s: missing %q", path, want)
				}
			}
			if strings.Contains(strings.ToLower(body), "tacyborg") || strings.Contains(body, "portfolio") {
				t.Errorf("%s must not serve the artist portfolio", path)
			}
		}
	}
}

func TestRoutes(t *testing.T) {
	handler := routes()
	for path, status := range map[string]int{
		"/static/site.css": 200, "/static/htmx.js": 200,
		"/en": 404, "/missing": 404, "/static/portfolio/web/example.webp": 404,
	} {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != status {
			t.Errorf("%s: got %d, want %d", path, w.Code, status)
		}
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /: got %d, want 405", w.Code)
	}
}
