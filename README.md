# CyborgNerd + Tacy Borg

Two independent websites in one repository, hosted on the same AWS Lightsail VM.

| App | Domain | Local address | Stack |
| --- | --- | --- | --- |
| [`tacyborg/`](tacyborg/) | `www.tacyborg.com` | `127.0.0.1:8080` | Go, HTML templates, portfolio CSS/JS |
| [`cyborgnerd/`](cyborgnerd/) | `www.cyborgnerd.com` | `127.0.0.1:8081` | Go, plain HTML/CSS + HTMX |

The artist portfolio previously under `src/` is now under `tacyborg/`, with its
content and assets preserved. The new CyborgNerd app serves only the personal
site; it has no portfolio routes or artwork assets.

## Run locally

Go 1.23.2+ is required. In separate terminals:

```sh
make -C tacyborg dev
```

```sh
make -C cyborgnerd dev
```

Both can run at once. `HTTP_ADDR` overrides an app's bind address. Templates and
static assets are embedded, so restart after editing. Neither app needs a frontend
build or npm installation. Tacy Borg's older gallery styles can optionally be
regenerated with its existing Tailwind commands.

## Personal site

Edit `cyborgnerd/templates/page.html` for content and `cyborgnerd/static/site.css`
for styling. Home, About and Projects are deliberately minimal starter pages;
replace the placeholder copy with your own bio and projects.

HTMX is served locally at `/static/htmx.js` (the existing bundled 2.0.3 version).
It enhances ordinary links by replacing the main content and updating browser
history. Every page is also directly accessible and works without JavaScript.
There is no custom application JavaScript, framework or Tailwind dependency.

## Test and build

From the repository root:

```sh
make test
GOOS=linux GOARCH=amd64 make binary
```

Outputs: `bin/tacyborg` and `bin/cyborgnerd`. Each app is a separate Go module;
use `make test` at the root or `go test ./...` inside an app directory.

With both local servers running and Playwright installed, optional browser checks:

```sh
NODE_PATH=/tmp/opencode/portfolio-browser/node_modules node tacyborg/tools/check_portfolio.cjs
NODE_PATH=/tmp/opencode/portfolio-browser/node_modules node cyborgnerd/tools/check_site.cjs
```

See `tacyborg/PORTFOLIO.md` for Playwright installation. Override `PORTFOLIO_URL`
or `CYBORGNERD_URL` to check a different local address.

## Deployment

Nginx routes the domains to separate loopback ports. Each app has its own systemd
service, OS user and install directory. Updates and restarts are independent.
See [deploy/README.md](deploy/README.md) for the **one-time production migration**,
backups, checks and rollback. Repository setup alone does not change production.

## Private assets

`source-assets/`, original PDFs, purchase ZIPs and binaries are ignored. Keep
private source material backed up separately. Runtime portfolio assets remain
inside `tacyborg/static/`. See [tacyborg/PORTFOLIO.md](tacyborg/PORTFOLIO.md) for
portfolio content, fonts, regeneration and browser checks.
