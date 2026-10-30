# Two sites, one AWS server

The repository is split; production still needs a one-time migration from the
old shared portfolio service. Preparing these files does **not** migrate the server.

| Site | App / binary | Service / user | Loopback upstream |
| --- | --- | --- | --- |
| `tacyborg.com`, `www.tacyborg.com` | `tacyborg/`, `bin/tacyborg` | `tacyborg` | `127.0.0.1:8080` |
| `cyborgnerd.com`, `www.cyborgnerd.com` | `cyborgnerd/`, `bin/cyborgnerd` | `cyborgnerd` | `127.0.0.1:8081` |

Both apps embed their own HTML and static files. Neither needs Node, Go, a database
or a container on the server. Nginx selects the app by hostname and handles HTTPS.
Only ports 80/443 are public; never expose 8080/8081 in the AWS firewall.

## Existing instance

- Lightsail `cyborgnerd`, Ubuntu 24.04, x86_64, static IP `34.235.217.80`.
- AWS profile `cyborg`, region `us-east-1`.
- SSH key `~/.ssh/cyborgnerd-lightsail.pem`.
- Known-host file `~/.ssh/cyborgnerd-lightsail-known_hosts`.
- Active TLS sites: `/etc/nginx/sites-available/tacyborg` and
  `/etc/nginx/sites-available/cyborgnerd`.
- Previously both domains used `cyborgnerd.service` on port 8080.

Keep the existing certificates, DNS, firewall restrictions and Certbot renewal.
`nginx.conf` (Tacy Borg) and `cyborgnerd.nginx.conf` are HTTP-only bootstrap
templates for **fresh** installations, not replacements for active TLS files.

## Build and upload

From the repository root:

```sh
make test
GOOS=linux GOARCH=amd64 make binary
scp -i ~/.ssh/cyborgnerd-lightsail.pem \
  -o UserKnownHostsFile=~/.ssh/cyborgnerd-lightsail-known_hosts \
  bin/tacyborg bin/cyborgnerd deploy/tacyborg.service deploy/cyborgnerd.service \
  ubuntu@34.235.217.80:/tmp/
```

Upload only binaries and service units, not source PDFs, artwork originals or
font purchase archives. The two binaries have distinct names; the personal-site
binary must **not** replace the live portfolio until it has its own service.

## One-time migration (on the existing server)

Schedule a brief interruption while transferring port 8080 to the portfolio's
new service. Inspect the active units and both Nginx sites before making changes.

### 1. Preserve the current working deployment

```sh
backup="/var/backups/cyborgnerd-split-$(date -u +%Y%m%dT%H%M%SZ)"
sudo mkdir -m 700 "$backup"
sudo cp -p /opt/cyborgnerd/cyborgnerd "$backup/portfolio-binary"
sudo cp -p /etc/systemd/system/cyborgnerd.service "$backup/cyborgnerd.service"
sudo cp -p /etc/nginx/sites-available/tacyborg "$backup/tacyborg.nginx"
sudo cp -p /etc/nginx/sites-available/cyborgnerd "$backup/cyborgnerd.nginx"
echo "Rollback directory: $backup"
```

Record that directory before continuing. Do not reuse or overwrite an old backup.

### 2. Give Tacy Borg its own service

If `tacyborg` already exists, inspect it rather than recreating it:

```sh
getent passwd tacyborg
# Only if the user does not exist:
sudo useradd --system --user-group --home-dir /opt/tacyborg --no-create-home --shell /usr/sbin/nologin tacyborg
```

If `/opt/tacyborg/tacyborg` or `tacyborg.service` already exists, back it up too
before installing. Otherwise:

```sh
sudo install -d -o root -g root -m 755 /opt/tacyborg
sudo install -o root -g root -m 755 /tmp/tacyborg /opt/tacyborg/tacyborg.next
sudo mv /opt/tacyborg/tacyborg.next /opt/tacyborg/tacyborg
sudo install -o root -g root -m 644 /tmp/tacyborg.service /etc/systemd/system/tacyborg.service
sudo systemctl daemon-reload
sudo systemctl stop cyborgnerd
sudo systemctl enable --now tacyborg
curl --fail --retry 5 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/ -o /dev/null
curl --fail https://www.tacyborg.com/en -o /dev/null
```

The Tacy Borg Nginx site stays on port 8080, so it needs no change. If these checks
fail, stop here and roll back; do not proceed to the personal-site switch.

### 3. Install the personal site on port 8081

The existing `cyborgnerd` service user and directory are reused. The old binary
and unit must already be backed up by step 1.

```sh
sudo install -o root -g root -m 755 /tmp/cyborgnerd /opt/cyborgnerd/cyborgnerd.next
sudo mv /opt/cyborgnerd/cyborgnerd.next /opt/cyborgnerd/cyborgnerd
sudo install -o root -g root -m 644 /tmp/cyborgnerd.service /etc/systemd/system/cyborgnerd.service
sudo systemctl daemon-reload
sudo systemctl enable --now cyborgnerd
curl --fail --retry 5 --retry-connrefused --retry-delay 1 http://127.0.0.1:8081/ -o /dev/null
```

Inspect `/etc/nginx/sites-available/cyborgnerd`. Change **every** app upstream in
that virtual host (including any separate `/static/` location) from
`127.0.0.1:8080` to `127.0.0.1:8081`. Preserve server names, certificates, redirects,
TLS directives and the separate Tacy Borg site.

```sh
sudoedit /etc/nginx/sites-available/cyborgnerd
sudo nginx -t
sudo systemctl reload nginx
```

### 4. Verify isolation

- Check HTTPS for all four hostnames, plus `/en` on Tacy Borg.
- Tacy Borg must still show the portfolio, images and documentary.
- CyborgNerd must show the personal site; `/about` and `/projects` must work both
  through HTMX navigation and with JavaScript disabled.
- CyborgNerd `/en` and `/static/portfolio/...` must return 404, not portfolio content.
- Check `systemctl is-active tacyborg cyborgnerd` and both service journals.
- Both services must remain enabled after reboot. Do not reboot just to test this.

### Roll back the migration

Use the actual backup path from step 1. This restores the old shared portfolio
on both domains, not the split setup:

```sh
sudo systemctl stop cyborgnerd
sudo systemctl disable --now tacyborg
sudo cp -p "$backup/portfolio-binary" /opt/cyborgnerd/cyborgnerd.next
sudo mv /opt/cyborgnerd/cyborgnerd.next /opt/cyborgnerd/cyborgnerd
sudo cp -p "$backup/cyborgnerd.service" /etc/systemd/system/cyborgnerd.service
sudo cp -p "$backup/tacyborg.nginx" /etc/nginx/sites-available/tacyborg
sudo cp -p "$backup/cyborgnerd.nginx" /etc/nginx/sites-available/cyborgnerd
sudo systemctl daemon-reload
sudo systemctl restart cyborgnerd
sudo nginx -t
sudo systemctl reload nginx
curl --fail --retry 5 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/ -o /dev/null
```

## Later updates

Build/upload only the app being changed. Save its current binary to a unique
backup, install the new binary to `<app>.next`, then rename it over the running
binary and restart **only that app's service**. Check loopback and public HTTPS.
If checks fail, restore the backup via the same rename process and restart.
No Nginx or certificate changes are needed for ordinary app updates.

Historical hashes and older rollback paths are in
[portfolio-deployment-history.md](portfolio-deployment-history.md).
