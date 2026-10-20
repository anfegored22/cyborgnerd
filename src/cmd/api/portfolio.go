package api

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
)

//go:embed portfolio.json
var portfolioSource []byte

//go:embed portfolio-layout.json
var portfolioLayoutSource []byte

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

func loadPortfolio() ([]portfolioSlide, error) {
	var slides []portfolioSlide
	if err := json.Unmarshal(portfolioSource, &slides); err != nil {
		return nil, err
	}
	var pages []portfolioPage
	if err := json.Unmarshal(portfolioLayoutSource, &pages); err != nil {
		return nil, err
	}
	if len(slides) != len(pages) {
		return nil, fmt.Errorf("portfolio content and PDF layout have different page counts")
	}
	for i := range slides {
		if slides[i].ID != pages[i].ID {
			return nil, fmt.Errorf("portfolio slide %d does not match PDF layout", i+1)
		}
		slides[i].Page = pages[i]
	}
	return slides, nil
}
