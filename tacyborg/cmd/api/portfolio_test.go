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
	"testing/fstest"
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
			if image.Width <= 0 || image.Height <= 0 {
				t.Errorf("missing stable image dimensions: %s", image.File)
			}
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
	if slides[4].ID != "proyectos" || slides[len(slides)-4].ID != "curriculum" || slides[len(slides)-3].ID != "exposiciones" || slides[len(slides)-2].ID != "experiencia" {
		t.Error("projects must follow research, with the résumé sections together before contact")
	}
}

func TestGenealogiaPhotographs(t *testing.T) {
	works := []struct {
		file, title, dimensions, year string
		width, height                 int
	}{
		{"genealogia-espejo", "Espejo", "100 cm × 86 cm", "2008", 1379, 1600},
		{"genealogia-leticia", "Leticia", "60 cm × 50 cm", "2008", 1600, 1323},
		{"genealogia-sin-titulo-2008", "Sin título", "60 cm × 50 cm", "2008", 1600, 1304},
		{"genealogia-la-rosa", "La Rosa", "24 cm × 14 cm", "2008", 1600, 994},
		{"genealogia-sin-titulo-2009", "Sin título", "100 cm × 70 cm", "2009", 1120, 1600},
		{"genealogia-abuela", "Abuela", "25 cm × 17 cm", "2009", 1069, 1600},
		{"genealogia-abuela-no-sufras", "Abuela, no sufras", "25 cm × 17 cm", "2009", 1075, 1600},
		{"genealogia-sin-nombre", "Sin nombre", "48 cm × 33 cm", "2007", 1600, 1079},
		{"genealogia-de-cama", "De cama", "47 cm × 30 cm", "2008", 1600, 1077},
		{"genealogia-el-barrio", "El barrio", "48 cm × 29 cm", "2008", 969, 1600},
	}
	for _, language := range []string{"es", "en"} {
		t.Run(language, func(t *testing.T) {
			slides, err := loadPortfolioLanguage(language)
			if err != nil {
				t.Fatal(err)
			}
			var images []portfolioImage
			for _, slide := range slides {
				if slide.ID == "genealogia" {
					images = slide.Images
				}
			}
			if len(images) != len(works)+1 {
				t.Fatalf("got %d images, want %d photographs", len(images), len(works)+1)
			}
			path := "/"
			if language == "en" {
				path = "/en"
			}
			w := httptest.NewRecorder()
			if err := portfolioTestServer(t).Hello(w, httptest.NewRequest(http.MethodGet, path, nil)); err != nil {
				t.Fatal(err)
			}
			if strings.Contains(w.Body.String(), "page-17-image-090.webp") {
				t.Error("repeated photograph must not appear on the rendered page")
			}
			for i, work := range works {
				image := images[i+1]
				technique := "Fotografía a blanco y negro"
				if language == "en" {
					technique = "Black-and-white photography"
				}
				if work.file == "genealogia-sin-titulo-2009" {
					if language == "en" {
						technique += " with color"
					} else {
						technique += " con color"
					}
				}
				caption := work.title + " · " + technique + " · " + work.dimensions + " · " + work.year
				if image.File != work.file || image.Caption != caption || image.Width != work.width || image.Height != work.height || image.Alt == "" {
					t.Errorf("incorrect photograph metadata: %+v", image)
				}
				if !strings.Contains(w.Body.String(), "/static/portfolio/web/"+work.file+".webp") || !strings.Contains(w.Body.String(), "<figcaption>"+caption+"</figcaption>") {
					t.Errorf("photograph or caption missing from rendered page: %s", work.file)
				}
				if info, err := os.Stat(filepath.Join("../../static/portfolio/web", work.file+".webp")); err != nil || info.Size() == 0 {
					t.Errorf("missing photograph asset: %s", work.file)
				}
			}
		})
	}
}

func TestEnglishPortfolio(t *testing.T) {
	s := portfolioTestServer(t)
	w := httptest.NewRecorder()
	if err := s.Hello(w, httptest.NewRequest(http.MethodGet, "/en", nil)); err != nil {
		t.Fatal(err)
	}
	body := w.Body.String()
	for _, expected := range []string{`lang="en"`, "Artist Statement", "Research Lines", "Project Portfolio", "Work Experience", "Portfolio options", "Read in Spanish", `href="https://www.tacyborg.com/en"`, "Watch the documentary on YouTube"} {
		if !strings.Contains(body, expected) {
			t.Errorf("English page is missing %q", expected)
		}
	}
	if w.Header().Get("Content-Language") != "en" || strings.Count(body, "data-slide data-title=") != 18 {
		t.Error("English page must contain all 18 sections with an English language header")
	}
	if strings.Contains(body, "ZgotmplZ") || strings.Contains(body, "<no value>") {
		t.Error("invalid English template data")
	}
	spanish, err := loadPortfolio()
	if err != nil {
		t.Fatal(err)
	}
	english, err := loadPortfolioLanguage("en")
	if err != nil {
		t.Fatal(err)
	}
	for i, slide := range english {
		if slide.ID != spanish[i].ID || slide.Page.ID != slide.ID || slide.Page.Background != spanish[i].Page.Background || len(slide.Page.Runs) == 0 {
			t.Errorf("English section/layout mismatch: %s", slide.ID)
		}
		for _, run := range slide.Page.Runs {
			if run.Size <= 0 || run.Left < 0 || run.Left >= 100 || run.Baseline < 0 || run.Baseline > 100 {
				t.Errorf("invalid English text geometry: %s %+v", slide.ID, run)
			}
		}
	}
	redirect := httptest.NewRecorder()
	if err := s.Hello(redirect, httptest.NewRequest(http.MethodGet, "/en/?view=test", nil)); err != nil {
		t.Fatal(err)
	}
	if redirect.Code != http.StatusPermanentRedirect || redirect.Header().Get("Location") != "/en?view=test" {
		t.Error("/en/ must redirect to /en and preserve query parameters")
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
	// Continuous scrolling keeps every page in the document, including without JS.
	if strings.Contains(body, "data-slide hidden") {
		t.Error("slides must remain available without JavaScript")
	}
	if index := strings.Index(body, "ZgotmplZ"); index >= 0 {
		t.Errorf("HTML template rejected generated layout style near: %.120s", body[max(0, index-70):])
	}
	if strings.Contains(body, "site-header") {
		t.Error("the original PDF view must not add a site header")
	}
	if strings.Contains(body, "portfolio-controls") || strings.Contains(body, "previous-slide") || strings.Contains(body, "next-slide") {
		t.Error("the continuous portfolio must not render the old bottom navigation")
	}
	if !strings.Contains(body, "portfolio-menu") {
		t.Error("missing floating section menu")
	}
	if !strings.Contains(body, `href="/static/portfolio.css?v=`) || !strings.Contains(body, `src="/static/js/portfolio.js?v=`) {
		t.Error("CSS and JavaScript need content versions to avoid mixing old and new navigation after deployment")
	}
	for _, text := range []string{"pdf-frame", "pdf-text-layer", "/static/portfolio/pages/", "toggle-reading", "zoom-slide", "artwork-link", `href="/static/portfolio/web/`} {
		if strings.Contains(body, text) {
			t.Errorf("obsolete PDF layout or image-opening link: %s", text)
		}
	}
	for _, text := range []string{`class="reading-view"`, `https://www.youtube-nocookie.com/embed/G9VE13QZeps?start=15`, `title="Documental de Historia humana"`, `allowfullscreen`, `https://www.youtube.com/watch?v=G9VE13QZeps&amp;t=15s`} {
		if !strings.Contains(body, text) {
			t.Errorf("missing responsive layout or documentary player: %s", text)
		}
	}
	if strings.Count(body, "<iframe ") != 1 || strings.Contains(body, "autoplay") {
		t.Error("only the documentary should have a player, without autoplay")
	}
}

func TestVersionedAssetURL(t *testing.T) {
	files := fstest.MapFS{"portfolio.css": {Data: []byte("first version")}}
	first, err := versionedAssetURL(files, "portfolio.css")
	if err != nil {
		t.Fatal(err)
	}
	files["portfolio.css"].Data = []byte("second version")
	second, err := versionedAssetURL(files, "portfolio.css")
	if err != nil {
		t.Fatal(err)
	}
	if first == second || !strings.HasPrefix(second, "/static/portfolio.css?v=") {
		t.Error("changing an asset must change its browser cache key")
	}
	if _, err := versionedAssetURL(files, "missing.css"); err == nil {
		t.Error("missing assets must return an error")
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
