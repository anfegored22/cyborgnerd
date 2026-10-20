# Tacy Borg portfolio

Run `go run .` from `src/`, then visit <http://localhost:8080>.
The new home page has its own CSS and does not require the Tailwind build.
The older gallery routes remain available and still use Tailwind.

## Content

- `cmd/api/portfolio.json`: editable Spanish content, in the PDF's 18-page order.
- `cmd/api/portfolio-layout.json`: generated original text positions, colors and sizes.
- `templates/home.html`: semantic HTML slides; all content works without JavaScript.
- `static/portfolio.css`: responsive layout, print styles and reduced-motion support.
- `static/js/portfolio.js`: buttons, section picker, horizontal swipe and keyboard navigation.
- `static/portfolio/text.txt`: original layout-preserving text extraction.
- `static/portfolio/manifest.json`: original image/page mapping.
- `static/portfolio/web/`: optimized WebP images used by the site.
- `static/portfolio/pages/`: text-free PDF backdrops with the exact artwork,
  gradients and vector decorations. Visible text is positioned HTML, not a screenshot.
- `../source-assets/portfolio/`: ignored original image extractions, kept outside
  the public/embedded files. Back these up separately; they are not in Git.
- `../source-assets/portfolio/unused-web/`: preserved derivatives not used by the site.

Each slide has a stable URL fragment, for example `/#chakana`. Arrow keys change
slides; Home/End jump to the first/last slide. Browser back/forward is supported.
On desktop, a plain mouse click on the original artwork or page background advances
one slide. Text, links, drags, modified clicks, and touch taps do not trigger this shortcut.
Click-to-advance is disabled in reading and zoomed modes so reading/panning is safe.
There is no autoplay. The default view preserves the original 16:9 page composition,
including headings that appear at the bottom and the PDF's original columns.
**Ampliar** shows the page at its original 1440px width, with horizontal/vertical
panning (swipe navigation is disabled while panning). **Lectura** switches to the
responsive, reflowed HTML version for small screens. Both modes have selectable text.
Printing includes all original-layout slides.

The original view uses text directly extracted from the PDF. Only the optional
reading view reorganizes text into readable paragraphs and columns.
Dates and wording from the PDF are preserved, including source inconsistencies
such as “Hiska 1”/“Hischa 1” and the Varietale date “2026–2017”; review these before publishing.

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

Only the 18 referenced reading-view images and 18 original-page backdrops need
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
