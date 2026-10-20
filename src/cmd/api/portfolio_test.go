package api

import (
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func portfolioTestServer(t *testing.T) *Server {
	t.Helper()
	tmpl, err := template.ParseGlob("../../templates/*.html")
	if err != nil {
		t.Fatal(err)
	}
	return &Server{Tmpl: tmpl, StaticFS: os.DirFS("../../static")}
}

func TestPortfolioContent(t *testing.T) {
	slides, err := loadPortfolio()
	if err != nil {
		t.Fatal(err)
	}
	if len(slides) != 18 {
		t.Fatalf("got %d slides, want 18 PDF pages", len(slides))
	}
	ids := make(map[string]bool)
	for _, slide := range slides {
		if slide.ID == "" || slide.Title == "" || slide.Kind == "" {
			t.Errorf("slide is missing required content: %+v", slide)
		}
		if ids[slide.ID] {
			t.Errorf("duplicate slide ID: %s", slide.ID)
		}
		ids[slide.ID] = true
		if slide.Page.ID != slide.ID || slide.Page.Width != 1440 || slide.Page.Height != 810 {
			t.Errorf("invalid PDF page geometry for %s", slide.ID)
		}
		if len(slide.Page.Runs) == 0 {
			t.Errorf("missing selectable PDF text for %s", slide.ID)
		}
		if slide.Page.Background != filepath.Base(slide.Page.Background) {
			t.Errorf("invalid PDF background filename: %s", slide.Page.Background)
		}
		background := filepath.Join("../../static/portfolio/pages", slide.Page.Background)
		if info, err := os.Stat(background); err != nil || info.Size() == 0 {
			t.Errorf("missing or empty PDF backdrop: %s", background)
		}
		for _, run := range slide.Page.Runs {
			if run.Text == "" || run.Size <= 0 || run.Width <= 0 || run.Fit <= 0 {
				t.Errorf("invalid text geometry on %s: %+v", slide.ID, run)
			}
			if run.Font != "pixelion" && run.Font != "questrial" {
				t.Errorf("unexpected font on %s: %s", slide.ID, run.Font)
			}
			if run.Left < 0 || run.Left >= 100 || run.Baseline < 0 || run.Baseline > 100 {
				t.Errorf("text origin outside page %s: %+v", slide.ID, run)
			}
		}
		for _, image := range slide.Images {
			if image.Alt == "" {
				t.Errorf("missing alt text: %s", image.File)
			}
			if image.File != filepath.Base(image.File) {
				t.Errorf("invalid image filename: %s", image.File)
			}
			path := filepath.Join("../../static/portfolio/web", image.File+".webp")
			if info, err := os.Stat(path); err != nil || info.Size() == 0 {
				t.Errorf("missing or empty image: %s", path)
			}
		}
	}
	if slides[0].ID != "portafolio" || slides[len(slides)-1].ID != "contacto" {
		t.Error("portfolio must start with the cover and end with contact information")
	}
}

func TestPortfolioHome(t *testing.T) {
	s := portfolioTestServer(t)
	w := httptest.NewRecorder()
	if err := s.Hello(w, httptest.NewRequest(http.MethodGet, "/", nil)); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("got status %d", w.Code)
	}
	body := w.Body.String()
	for _, text := range []string{`lang="es"`, "Tacy Borg", "Chakana", "Cxacxa", "Genealogía del silencio", "mailto:tacyborg@outlook.com", "/static/portfolio.css", "/static/js/portfolio.js"} {
		if !strings.Contains(body, text) {
			t.Errorf("home is missing %q", text)
		}
	}
	if count := strings.Count(body, "data-slide data-title="); count != 18 {
		t.Errorf("rendered %d slides, want 18", count)
	}
	// All content is present without JavaScript; enhancement hides inactive slides.
	if strings.Contains(body, "data-slide hidden") {
		t.Error("slides must remain available without JavaScript")
	}
	if index := strings.Index(body, "ZgotmplZ"); index >= 0 {
		t.Errorf("HTML template rejected generated layout style near: %.120s", body[max(0, index-70):])
	}
	if strings.Contains(body, "site-header") {
		t.Error("the original PDF view must not add a site header")
	}
	for _, text := range []string{"pdf-text-layer", "pdf-run font-questrial", "pdf-run font-pixelion", "toggle-reading", "zoom-slide"} {
		if !strings.Contains(body, text) {
			t.Errorf("missing selectable layout or controls: %s", text)
		}
	}
}

func TestPortfolioAssets(t *testing.T) {
	s := portfolioTestServer(t)
	routes := s.routes()
	for _, path := range []string{
		"/static/portfolio.css",
		"/static/js/portfolio.js",
		"/static/fonts/Questrial-Regular.ttf",
		"/static/fonts/Silkscreen-Regular.ttf",
		"/static/fonts/Pixelion-Regular.woff2",
		"/static/fonts/Pixelion-Regular.woff",
		"/static/portfolio/text.txt",
		"/static/portfolio/web/page-10-image-052.webp",
		"/static/portfolio/pages/page-01.webp",
		"/static/portfolio/pages/page-18.webp",
	} {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			routes.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
			if w.Code != http.StatusOK {
				t.Errorf("got %d for %s", w.Code, path)
			}
		})
	}
	if _, err := fs.Stat(s.StaticFS, "fonts/questrial-OFL.txt"); err != nil {
		t.Error("missing Questrial license:", err)
	}
	if _, err := fs.Stat(s.StaticFS, "fonts/silkscreen-OFL.txt"); err != nil {
		t.Error("missing Silkscreen license:", err)
	}
}

func TestUnknownRoute(t *testing.T) {
	s := portfolioTestServer(t)
	w := httptest.NewRecorder()
	if err := s.Hello(w, httptest.NewRequest(http.MethodGet, "/missing-page", nil)); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusNotFound {
		t.Errorf("got status %d, want 404", w.Code)
	}
}
