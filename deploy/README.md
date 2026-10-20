# Instance deployment: Go + systemd + Nginx

No containers, database, or persistent application disk are required. The executable
embeds the templates, optimized images and fonts. Nginx proxies to Go on loopback;
systemd runs it as an unprivileged user and restarts it on failure.

These instructions target Ubuntu/Debian with systemd. They prepare one hostname;
replace `portfolio.example.com` with yours. No commands below have been run on a
remote machine by the coding agent.

## 1. Build and upload (your computer)

From `src/`, with Go 1.23.2 or newer installed:

```sh
go test ./...
GOOS=linux GOARCH=amd64 make binary
```

Use `GOARCH=arm64` instead for an ARM instance. Check `uname -m` on the server:
`x86_64` means amd64, `aarch64` means arm64. The output is `bin/cyborgnerd` at the
repository root, built with CGO disabled. It does not need Go installed to run.

The portfolio uses plain CSS. If you also want the old `/artist-path`, `/assets`
and gallery pages styled, generate their Tailwind CSS with `make css` **before**
building the binary. That optional command needs the existing Tailwind v3 CLI and
the forms/typography plugins; it is not needed for the portfolio itself.

From the repository root, upload the binary and configs using your SSH user/host:

```sh
scp bin/cyborgnerd deploy/cyborgnerd.service deploy/nginx.conf USER@SERVER:/tmp/
```

Upload only these deployment files, not the repository, original artwork, source
PDF or purchased font archive. If a service/config file already exists on your
server, inspect and back it up before replacing it.

## 2. Install the service (on the server)

```sh
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx
id cyborgnerd
```

If the user does not exist, create it:

```sh
sudo useradd --system --user-group --home-dir /opt/cyborgnerd --no-create-home --shell /usr/sbin/nologin cyborgnerd
```

Install the executable as root-owned and readable/executable by the service user:

```sh
sudo install -d -o root -g root -m 755 /opt/cyborgnerd
sudo install -o root -g root -m 755 /tmp/cyborgnerd /opt/cyborgnerd/cyborgnerd
sudo install -o root -g root -m 644 /tmp/cyborgnerd.service /etc/systemd/system/cyborgnerd.service
sudo systemctl daemon-reload
sudo systemctl enable --now cyborgnerd
curl --fail --retry 5 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/ -o /dev/null
```

Go listens on `127.0.0.1:8080`; do not expose port 8080 in your cloud firewall.
For local development, `HTTP_ADDR` can override this address.

## 3. Configure Nginx and HTTPS

Point your domain's DNS A record to the instance. Set an AAAA record only if the
instance has working public IPv6. Allow inbound HTTP/HTTPS (80/443) and restrict
SSH access appropriately in the instance firewall/security group. If you enable
UFW, allow SSH before enabling it so you do not lock yourself out.

```sh
sudo install -o root -g root -m 644 /tmp/nginx.conf /etc/nginx/sites-available/cyborgnerd
sudoedit /etc/nginx/sites-available/cyborgnerd
```

Replace `portfolio.example.com` with your real domain. Then enable this site:

```sh
sudo ln -s /etc/nginx/sites-available/cyborgnerd /etc/nginx/sites-enabled/cyborgnerd
sudo nginx -t
sudo systemctl enable --now nginx
sudo systemctl reload nginx
sudo certbot --nginx -d YOUR_DOMAIN --redirect
sudo certbot renew --dry-run
curl --fail https://YOUR_DOMAIN/ -o /dev/null
```

Do not remove other Nginx sites unless you have confirmed they are unused. If the
symlink already exists, skip creating it. Certbot needs DNS to point to this server
and port 80 reachable before it can issue the certificate. Check that the Certbot
renewal timer is enabled (`systemctl list-timers --all | grep certbot`).

## 4. Update and roll back

Build/test locally and upload a new binary to `/tmp/cyborgnerd`. On the server:

```sh
sudo cp /opt/cyborgnerd/cyborgnerd /opt/cyborgnerd/cyborgnerd.previous
sudo install -o root -g root -m 755 /tmp/cyborgnerd /opt/cyborgnerd/cyborgnerd.next
sudo mv /opt/cyborgnerd/cyborgnerd.next /opt/cyborgnerd/cyborgnerd
sudo systemctl restart cyborgnerd
curl --fail --retry 5 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/ -o /dev/null
curl --fail https://YOUR_DOMAIN/ -o /dev/null
```

Replacing the executable via rename avoids overwriting the file used by the
running process. If verification fails, restore the saved version:

```sh
sudo cp /opt/cyborgnerd/cyborgnerd.previous /opt/cyborgnerd/cyborgnerd.next
sudo mv /opt/cyborgnerd/cyborgnerd.next /opt/cyborgnerd/cyborgnerd
sudo systemctl restart cyborgnerd
```

This single-instance restart has brief downtime. No Nginx reload is required for
binary-only updates. Static assets are cached for up to one hour; hard-refresh the
browser when checking a new release.

## Logs and troubleshooting

```sh
sudo systemctl status cyborgnerd
sudo journalctl -u cyborgnerd -f
sudo tail -f /var/log/nginx/error.log
```

Keep the ignored `source-assets/`, source PDF and purchase/license documents backed
up separately. Git contains the optimized runtime images and regeneration scripts,
not those originals or compiled binaries. Deployment does not require originals.
