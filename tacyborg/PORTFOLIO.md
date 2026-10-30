# Tacy Borg portfolio

## Current presentation

All screen sizes now use the responsive reading layout, with selectable HTML
text and individual artwork images. The PDF composition and reading/zoom switches
are no longer rendered. Scrolling is vertical, with the floating section menu,
deep links, keyboard section navigation and language switch retained. The same
layout is available without JavaScript and when printing.

Artwork images do not open on click or tap. The former
`page-15-image-070.webp` image is replaced by a lazy-loaded, privacy-enhanced
YouTube documentary player, starting at 15 seconds without autoplay. A regular
YouTube link remains underneath in case embedding is unavailable.

The PDF layout data, backdrops and generation tools below are retained as source
material, but are not requested or displayed by the home page. The historical
original-view behavior described below no longer applies to the live template.

Run `go run .` from `tacyborg/`, then visit <http://localhost:8080>.
The new home page has its own CSS and does not require the Tailwind build.
The older gallery routes remain available and still use Tailwind.

## Content

- `cmd/api/portfolio.json`: editable Spanish content, in the PDF's 18-page order.
- `cmd/api/portfolio.en.json`: English content, artwork descriptions and captions.
- `cmd/api/portfolio-layout.json`: generated original text positions, colors and sizes.
- `cmd/api/portfolio-layout.en.json`: English text fitted into the same artwork layouts.
- `cmd/api/portfolio_text.go`: Spanish/English interface labels and instructions.
- `templates/home.html`: semantic HTML slides; all content works without JavaScript.
- `static/portfolio.css`: responsive layout, print styles and reduced-motion support.
- `static/js/portfolio.js`: continuous scrolling, nearby-image decoding, section picker and keyboard navigation.
- `static/portfolio/text.txt`: original layout-preserving text extraction.
- `static/portfolio/manifest.json`: original image/page mapping.
- `static/portfolio/web/`: optimized WebP images used by the site.
- `static/portfolio/pages/`: text-free PDF backdrops with the exact artwork,
  gradients and vector decorations. Visible text is positioned HTML, not a screenshot.
- `../source-assets/portfolio/`: ignored original image extractions, kept outside
  the public/embedded files. Back these up separately; they are not in Git.
- `../source-assets/portfolio/unused-web/`: preserved derivatives not used by the site.
- `../source-assets/genealogia/`: private JPEG originals from `imagenesparapaginaweb.zip`,
  renamed to match the `genealogia-*.webp` derivatives. The seven added photographs
  follow the existing images in Genealogía del silencio, with technical captions
  in both languages; artwork titles remain in Spanish.

All 18 pages stay in one continuous document; there is no hidden-page slideshow,
click-to-advance, autoplay, looping, or bottom navigation bar. Wide viewports
(at least 900px and aspect ratio 4:3) scroll horizontally; narrower/portrait
viewports scroll vertically. A vertical mouse wheel maps to horizontal movement
in the wide original view, while horizontal trackpads and touch scrolling remain
native. Gentle snapping does not force every small wheel step back to a page edge.

Narrow/portrait screens automatically use the reflowed reading view, rather than
shrinking a landscape PDF into unreadably small text. Phones use 18px body text
and naturally sized, full-width artwork. Wide screens keep the PDF composition.
Resizing switches automatically unless the visitor explicitly chooses a view
with the Lectura toggle; that choice is then preserved for the current visit.

Each page has a stable URL fragment, for example `/#chakana`. Scrolling updates it
without adding a history entry for every page. Explicit section jumps preserve
browser back/forward. Arrow keys navigate pages, with Up/Down supported in the
vertical original view; Home/End jump to the first/last page. Reduced-motion users
get instant section jumps. Resizing preserves the current section.

A small floating **Opciones** menu replaces the footer. It contains the section
picker, **Ampliar** (1440px original pages with native panning), and **Lectura**
(responsive, reflowed text). Reading and zoom modes always use a vertical document;
reading has no nested scroll containers or page snapping to trap long text.
Both modes retain selectable text. Printing includes all original-layout pages.

The current page, previous page, and next two pages are requested and decoded
ahead of arrival. Only images for the active view are warmed; the first backdrop
is eager/high-priority and the rest start lazy. Fixed canvas dimensions and image
width/height metadata reserve space; decoded images fade in without layout jumps.
Very slow or failed requests do not block scrolling. Nearby prefetching cannot
guarantee a distant page is already loaded after an immediate section jump.

The CSS and JavaScript URLs include content hashes. A deployment that changes
the controls therefore fetches matching assets instead of reusing the old cached
footer-navigation code. Production must be explicitly updated after local review.

The original view uses text directly extracted from the PDF. Only the optional
reading view reorganizes text into readable paragraphs and columns.
Dates and wording from the PDF are preserved, including source inconsistencies
such as “Hiska 1”/“Hischa 1” and the Varietale date “2026–2017”; review these before publishing.

## Order and languages

The displayed order is defined once in `cmd/api/portfolio.go`: introduction and
research, project portfolio, then CV/exhibitions/work experience, ending with
contact. Both languages share that order. Source content/layout files retain
PDF order; they are joined by stable section IDs, not by array positions.

Spanish is `/`; English is `/en`. `/en/` redirects to `/en`. The visible EN/ES
switch carries the section fragment and reading/zoom choice across languages.
It also works without JavaScript (without preserving the view state). Translations
cover both the responsive reading view and the original-layout HTML text, plus
navigation, accessibility labels, image descriptions and metadata. Artwork titles
and institutional names are retained where they are proper names.

The canonical URLs use `https://www.tacyborg.com/` and `/en`, with alternate-language
links. All artwork and font assets are shared, not duplicated per language.

After editing English copy, regenerate its fitted original-layout text:

```sh
uv run tools/sync_english_layout.py
```

This needs the ignored source PDF and purchased `font files.zip`, or the script's
`--pdf` / `--font-zip` arguments. It only writes English text geometry; it reuses
the existing text-free backdrops and never extracts fonts from the PDF.

## Fonts

The PDF uses **Questrial Regular** and **Pixelion Regular**.
Questrial is bundled locally with its SIL Open Font License. Pixelion's WOFF2 and
WOFF files were supplied in `font files.zip` and are served locally from
`static/fonts/`; the archive is retained unchanged. These are separately supplied
webfonts, not extracted PDF subsets. Pixelion remains subject to its commercial
license; do not treat it as an open-source asset or redistribute the purchase archive.

The WOFF2 font is preloaded, with WOFF as a format fallback. The script detects
whether Pixelion loads and switches to its original PDF size and baseline.
**Silkscreen Regular**, bundled with its OFL license, remains a fitted substitute
if the Pixelion request fails. No third-party font requests are made by the site.

## Regenerating derivatives

With ImageMagick and Poppler installed:

```sh
python3 tools/prepare_portfolio.py
```

This script reads the extracted originals from `../source-assets/portfolio/` and
regenerates only the reading-view images referenced in `cmd/api/portfolio.json`,
at 1600px maximum, plus the raw text. Use `--originals PATH` if your originals are
stored elsewhere. It expects `../portfolio.pdf` unless `--pdf PATH` is supplied.
The raw PDF is not needed to run the website.

Additional photographs are read from `../source-assets/genealogia/` (override with
`--additional-originals PATH`); they are not part of the original PDF manifest.

Only the referenced reading-view images and 18 original-page backdrops need
to be deployed. The server's embed patterns exclude original PNG/JPG extractions.
Source PDFs, purchase ZIPs and compiled binaries are ignored as well. Keep originals
and license documents backed up privately; a fresh Git checkout does not contain them.

To regenerate the PDF-faithful view, install [uv](https://docs.astral.sh/uv/) and run:

```sh
uv run tools/sync_pdf_layout.py
```

Its Python dependencies are declared in the script. It reads the PDF without
modifying it, records text coordinates (not font files), and removes only text
before rendering the backdrops. All overlapping images and vector art are retained.
It writes `portfolio-layout.json` and the 18 backdrops. Regenerate after changing
the source PDF; `portfolio.json` is independently maintained for the reading view.

## Checks

```sh
go test ./cmd/api
node --check static/js/portfolio.js
```

Optional browser regression checks, with the Go server running locally:

```sh
npm install --prefix /tmp/opencode/portfolio-browser playwright
/tmp/opencode/portfolio-browser/node_modules/.bin/playwright install chromium
NODE_PATH=/tmp/opencode/portfolio-browser/node_modules node tools/check_portfolio.cjs
```

Set `PORTFOLIO_URL` to use another local port, and `PORTFOLIO_SCREENSHOTS` to an
output path prefix if you want desktop/mobile screenshots. The checks cover
desktop/mobile in both languages, vertical scrolling, section navigation and
deep links, resizing, non-clickable artwork, documentary markup, print and no-JS.
YouTube is stubbed in browser checks; verify real playback manually because the
video owner or YouTube may restrict embedding or regional availability.
