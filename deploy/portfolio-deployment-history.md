# Historical single-app deployment (before the repository split)

This records the previous deployment and rollback binaries. Paths and commands
below describe the old `src/` layout, not the new two-app setup. Do not use them
to deploy the split apps; follow [README.md](README.md) instead.

No containers, database, or persistent application disk are required. The executable
embeds the templates, optimized images and fonts. Nginx proxies to Go on loopback;
systemd runs it as an unprivileged user and restarts it on failure.

These instructions target Ubuntu/Debian with systemd. The bootstrap Nginx config
targets `tacyborg.com` and `www.tacyborg.com`. Use the initial installation
sections only for a fresh server; subsequent releases should use section 4 and
must not overwrite the active Certbot-managed TLS configuration.

## Lightsail instance

The selected VM is `cyborgnerd` in `us-east-1` (`us-east-1a`): Ubuntu 24.04 LTS,
512 MB RAM, 2 burstable vCPUs, and a 20 GB disk. Bundle `nano_3_0` costs $5/month
before taxes, optional backups, and transfer overages. No paid add-ons were enabled.

- Static IPv4: `34.235.217.80` (Lightsail resource `cyborgnerd-ip`, attached).
- Architecture verified over SSH: `x86_64`; build with `GOARCH=amd64`.
- AWS CLI profile: `cyborg`; always supply `--region us-east-1`.
- Dedicated SSH private key: `~/.ssh/cyborgnerd-lightsail.pem`, stored locally only.
- AWS key-pair name: `cyborgnerd-deploy`.

Connect using the host-key file recorded during the initial SSH verification:

```sh
ssh -i ~/.ssh/cyborgnerd-lightsail.pem -o UserKnownHostsFile=~/.ssh/cyborgnerd-lightsail-known_hosts ubuntu@34.235.217.80
```

SSH is allowed only from the workstation's public IPv4 at provisioning time. If
your public IP changes, update the port-22 source in the Lightsail firewall before
connecting; do not open SSH to the entire internet. Ports 80/443 accept IPv4/IPv6.
The private key and source-asset backups must not be committed to Git.

On this small instance, build the executable locally rather than installing the Go
compiler on the server. A 1 GB `/swapfile` has been enabled and added to `/etc/fstab`
for memory headroom during OS maintenance. It uses the included disk, not an extra
AWS volume. The original fstab is saved as `/etc/fstab.before-cyborgnerd`.

The domain uses Namecheap BasicDNS: `@` is an A record for `34.235.217.80`, and
`www` is a CNAME for `tacyborg.com`. No Route 53 hosted zone is required.

## Current deployment

The primary domain is `www.tacyborg.com`, with Spanish at `/` and English at `/en`.
Its Namecheap records point to the same instance. The old `cyborgnerd.com` and
`www.cyborgnerd.com` sites and certificates remain intact during migration.

- Live at <https://www.tacyborg.com> and <https://www.tacyborg.com/en>.
  `tacyborg.com` works too. HTTP redirects to HTTPS for both new hostnames.
- `/opt/cyborgnerd/cyborgnerd` runs as the unprivileged `cyborgnerd` user under
  systemd, listening only on `127.0.0.1:8080`. Nginx handles public 80/443.
- Both services and the Certbot renewal timer are enabled at boot.
  The new domain's simulated renewal passed. Keep port 80 reachable for HTTP
  certificate validation.
- The new certificate covers `tacyborg.com` and `www.tacyborg.com`, expires on
  January 19, 2027, and is stored under `/etc/letsencrypt/live/tacyborg.com/`;
  automatic renewal is configured, so no manual replacement should be necessary.
- Certbot was registered without an email address; add one later with
  `sudo certbot update_account --email YOUR_EMAIL` if desired.
- The initial release includes generated Tailwind CSS for the older galleries,
  in addition to the portfolio's standalone CSS. Generate it before rebuilding
  from a fresh checkout if those routes should remain styled.
- Initial deployed binary SHA-256:
  `c477ca7bf03a9d92b040277ff60269bdfade7fc82b7d7ec4ad6e0c8e1252da38`.
- Current release uses the responsive reading layout on all screens, non-clickable
  artwork, and an embedded documentary starting at 15 seconds, in Spanish and English.
  Deployed October 20, 2026; binary SHA-256:
  `e5b61f5e71d61ef356fba74efca97db7fee278cd20f69684760b75780aa09951`.
- The preceding bilingual release is preserved for rollback at
  `/opt/cyborgnerd/cyborgnerd.backup.Ue3cPR` (SHA-256
  `ae2f18f8b34c766655ff729deb2e2fe2574f2d65e39abd95ff7270cd6d27caab`).
  Public HTTPS and matching CSS/JS were verified on both new hostnames/languages;
  desktop/mobile browser checks passed against production. YouTube playback itself
  needs manual verification (the automated browser checks stub the external player).
- The previous mobile/scrolling release is preserved for rollback at
  `/opt/cyborgnerd/cyborgnerd.backup.RjxOnN` (SHA-256
  `8585af0bd097835bc8488418d4cecab55032fbb4a61526c0010f26e29eb3ea16`).
  The initial release is also retained at `/opt/cyborgnerd/cyborgnerd.backup.3856O4`.
  Public HTTPS, language headers, section order and matching versioned CSS/JS were
  verified for the new domain; the previous domain was checked too.

Active TLS configuration is `/etc/nginx/sites-available/tacyborg`, modified by
Certbot. The old site remains in `/etc/nginx/sites-available/cyborgnerd`.
`deploy/nginx.conf` is an HTTP-only bootstrap file, **not** a replacement
for the active TLS configuration. For normal updates, upload/restart the binary
only, using section 4.

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

For a first installation, upload the binary and configs from the repository root:

```sh
scp -i ~/.ssh/cyborgnerd-lightsail.pem \
  -o UserKnownHostsFile=~/.ssh/cyborgnerd-lightsail-known_hosts \
  bin/cyborgnerd deploy/cyborgnerd.service deploy/nginx.conf \
  ubuntu@34.235.217.80:/tmp/
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

Confirm the config names both `tacyborg.com` and `www.tacyborg.com`, then enable
this site:

```sh
sudo ln -s /etc/nginx/sites-available/cyborgnerd /etc/nginx/sites-enabled/cyborgnerd
sudo nginx -t
sudo systemctl enable --now nginx
sudo systemctl reload nginx
sudo certbot --nginx -d tacyborg.com -d www.tacyborg.com --redirect
sudo certbot renew --dry-run
curl --fail https://www.tacyborg.com/ -o /dev/null
```

Do not remove other Nginx sites unless you have confirmed they are unused. If the
symlink already exists, skip creating it. Certbot needs DNS to point to this server
and port 80 reachable before it can issue the certificate. Check that the Certbot
renewal timer is enabled (`systemctl list-timers --all | grep certbot`).

## Add SSL for the new domain on the existing server

Do not overwrite the working `cyborgnerd` Nginx site. Add a separate `tacyborg`
site so the existing domain/certificate remains available during migration.
If a `tacyborg` site already exists, inspect it first and do not overwrite an
existing TLS configuration. Skip symlink creation if it already exists.

Upload the new HTTP bootstrap configuration from the repository root:

```sh
scp -i ~/.ssh/cyborgnerd-lightsail.pem \
  -o UserKnownHostsFile=~/.ssh/cyborgnerd-lightsail-known_hosts \
  deploy/nginx.conf ubuntu@34.235.217.80:/tmp/tacyborg.nginx.conf
```

On the server:

```sh
sudo install -m 644 /tmp/tacyborg.nginx.conf /etc/nginx/sites-available/tacyborg
sudo ln -s /etc/nginx/sites-available/tacyborg /etc/nginx/sites-enabled/tacyborg
sudo nginx -t
sudo systemctl reload nginx
sudo certbot --nginx --cert-name tacyborg.com -d tacyborg.com -d www.tacyborg.com --redirect
sudo certbot renew --dry-run
curl --fail https://www.tacyborg.com/ -o /dev/null
```

This issues a separate free certificate covering both new hostnames. The existing
renewal timer handles it too. No certificate purchase in Namecheap is needed.
Deploying the new Go binary is a separate step; `/en` is available only after that
binary has been installed. Keep both old and new DNS records until migration is
verified. Redirecting the old domain can be done after the new site is working.

## 4. Update and roll back

Build/test locally, then upload the new binary to `/tmp/cyborgnerd`:

Upload from the repository root (do not upload the bootstrap Nginx config again):

```sh
scp -i ~/.ssh/cyborgnerd-lightsail.pem \
  -o UserKnownHostsFile=~/.ssh/cyborgnerd-lightsail-known_hosts \
  bin/cyborgnerd ubuntu@34.235.217.80:/tmp/cyborgnerd
```

Then run on the server:

```sh
sudo cp /opt/cyborgnerd/cyborgnerd /opt/cyborgnerd/cyborgnerd.previous
sudo install -o root -g root -m 755 /tmp/cyborgnerd /opt/cyborgnerd/cyborgnerd.next
sudo mv /opt/cyborgnerd/cyborgnerd.next /opt/cyborgnerd/cyborgnerd
sudo systemctl restart cyborgnerd
curl --fail --retry 5 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/ -o /dev/null
curl --fail https://www.tacyborg.com/ -o /dev/null
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
