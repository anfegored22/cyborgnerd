package api

import (
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
)

//go:embed portfolio.json
var portfolioSource []byte

//go:embed portfolio-layout.json
var portfolioLayoutSource []byte

//go:embed portfolio.en.json
var portfolioEnglishSource []byte

//go:embed portfolio-layout.en.json
var portfolioEnglishLayoutSource []byte

// Content files retain source-PDF order for regeneration. Presentation order
// is shared by both languages and joined to layouts by stable section ID.
var portfolioOrder = []string{
	"portafolio", "declaracion", "semblanza", "investigacion",
	"proyectos", "mosaicos", "chakana", "cxacxa", "rupestres", "hischa",
	"comunidad", "historia-humana", "amados-monstruos", "genealogia",
	"curriculum", "exposiciones", "experiencia", "contacto",
}

type portfolioTextRun struct {
	Text     string
	Font     string
	Left     float64
	Baseline float64
	Size     float64
	Width    float64
	Fit      float64
	Color    string
	// Generated contact links are trusted source data, including the tel scheme.
	Href template.URL
}

type portfolioPage struct {
	ID         string
	Background string
	Width      float64
	Height     float64
	Runs       []portfolioTextRun
}

type portfolioImage struct {
	File    string
	Alt     string
	Caption string
	Width   int
	Height  int
}

type portfolioColumn struct {
	Title      string
	Paragraphs []string
}

type portfolioSlide struct {
	ID         string
	Title      string
	Eyebrow    string
	Kind       string
	Paragraphs []string
	Columns    []portfolioColumn
	Images     []portfolioImage
	Quote      string
	Page       portfolioPage
}

type portfolioView struct {
	Slides    []portfolioSlide
	StyleURL  string
	ScriptURL string
	Language  string
	Canonical string
	OtherURL  string
	OtherLang string
	OtherCode string
	Text      map[string]string
}

func versionedAssetURL(files fs.FS, path string) (string, error) {
	content, err := fs.ReadFile(files, path)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(content)
	return fmt.Sprintf("/static/%s?v=%x", path, digest[:6]), nil
}

func loadPortfolio() ([]portfolioSlide, error) {
	return loadPortfolioLanguage("es")
}

func loadPortfolioLanguage(language string) ([]portfolioSlide, error) {
	source, layout := portfolioSource, portfolioLayoutSource
	if language == "en" {
		source, layout = portfolioEnglishSource, portfolioEnglishLayoutSource
	} else if language != "es" {
		return nil, fmt.Errorf("unsupported portfolio language: %s", language)
	}
	var slides []portfolioSlide
	if err := json.Unmarshal(source, &slides); err != nil {
		return nil, err
	}
	var pages []portfolioPage
	if err := json.Unmarshal(layout, &pages); err != nil {
		return nil, err
	}
	if len(slides) != len(pages) || len(slides) != len(portfolioOrder) {
		return nil, fmt.Errorf("portfolio content and PDF layout have different page counts")
	}
	byID := make(map[string]portfolioSlide, len(slides))
	for _, slide := range slides {
		if _, duplicate := byID[slide.ID]; duplicate {
			return nil, fmt.Errorf("duplicate portfolio section: %s", slide.ID)
		}
		byID[slide.ID] = slide
	}
	pageByID := make(map[string]portfolioPage, len(pages))
	for _, page := range pages {
		if _, duplicate := pageByID[page.ID]; duplicate {
			return nil, fmt.Errorf("duplicate portfolio page: %s", page.ID)
		}
		pageByID[page.ID] = page
	}
	ordered := make([]portfolioSlide, 0, len(slides))
	for _, id := range portfolioOrder {
		slide, found := byID[id]
		page, pageFound := pageByID[id]
		if !found || !pageFound {
			return nil, fmt.Errorf("missing portfolio section or layout: %s", id)
		}
		slide.Page = page
		ordered = append(ordered, slide)
	}
	return ordered, nil
}
