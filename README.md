# CyborgNerd

CyborgNerd is a Go web app featuring Tacy Borg's artistic portfolio, alongside the original art notes, sketches, and game assets.
The app serves server-rendered HTML templates with embedded static files and Tailwind CSS styling.

## Features

- Server-rendered pages using Go templates
- Embedded `templates/` and `static/` assets via `go:embed`
- Gallery-style routes for art notes and sketches
- Assets page for logos and animated sprites
- Simple error fallback page for unfinished sections

## Routes

- `/` - Tacy Borg portfolio carousel, with selectable text and original PDF layouts
- `/artist-path` - Artist landing page
- `/udemy-notes` - Udemy notes image gallery
- `/sketches` - Sketches image gallery
- `/assets` - Logos and animations
- `/dev-path` - Placeholder route (currently returns an under-construction page)

## Tech Stack

- Go 1.23
- Standard library `net/http` server
- Go templates (`html/template`)
- Tailwind CSS
- HTMX (served from local static JS)

## Project Layout

```text
.
├── src/
│   ├── cmd/api/           # HTTP server, handlers, route wiring
│   ├── templates/         # HTML templates
│   ├── static/            # CSS, JS, images, galleries
│   ├── main.go            # App entrypoint
│   └── Makefile           # Local dev/build commands
├── deploy/                # Nginx + systemd configs and instance deployment guide
├── bin/                   # Local compiled binaries (ignored)
└── source-assets/         # Private original artwork and unused derivatives (ignored)
```

## Prerequisites

- Go 1.23.2+
- Optional: Tailwind CSS v3 CLI and forms/typography plugins for the older gallery pages

## Run Locally

From `src/`, start the app:

```bash
make dev
```

Then open `http://localhost:8080`. The portfolio does not need a CSS build.
For the older gallery pages, run Tailwind in a second terminal:

```bash
make tw
```

The server binds to `127.0.0.1:8080` by default. Set `HTTP_ADDR` to override it.
Templates and assets are embedded: restart the app after changing them.

## Build

From `src/`:

```bash
make test
make binary
```

The binary is written to `bin/cyborgnerd` at the repository root. It embeds runtime
assets and does not require the original artwork, PDF, or font purchase archive.
`make build` additionally generates the optional legacy Tailwind CSS first.

## Deployment

Run the executable directly on a Linux instance, supervised by **systemd**, behind
**Nginx** with HTTPS from Certbot. No containers are used.

See [deploy/README.md](deploy/README.md) for installation, DNS/TLS setup, updates,
rollback, and logs. Configuration files are included under `deploy/`.

## Asset tracking

Commit the optimized WebP files, webfonts (subject to their license), templates,
content and tooling. Do not commit compiled executables, original extractions,
source PDFs or purchase archives. Keep ignored source assets backed up separately:
**being ignored by Git is not a backup**.

The extracted originals are preserved under `source-assets/portfolio/`, outside
the embedded public files. Unused WebP derivatives are in its `unused-web/` folder.
For asset regeneration and portfolio details, see [src/PORTFOLIO.md](src/PORTFOLIO.md).
