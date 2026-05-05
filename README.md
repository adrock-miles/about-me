# about-me

Personal portfolio site. Single Go binary, server-rendered HTML, no SPA, no database, no public API.

## Stack

- **Go 1.22** — chi router, `log/slog` JSON logs, Cobra + Viper for CLI and config
- **html/template** — parsed once at startup, embedded in the binary via `//go:embed`
- **Tailwind v4** browser CDN, paired with a native-CSS design system at `internal/web/static/css/styles.css`
- **Datastar** — all client interactivity (no React, no jQuery, no inline `<script>` for state)
- **Pluggable mailer** — `Sender` interface with Resend / SMTP / dev-mode log providers

The contract between browser and server is rendered HTML plus Datastar SSE patches. No JSON APIs.

## Layout

```
cmd/server/main.go     # cobra `serve` subcommand, the only entrypoint
internal/
  platform/            # config, server, middleware, mailer, sse, ratelimit
  handlers/            # HTTP handlers (page, contact, health)
  content/             # static project list
  web/                 # embedded templates + static assets
configs/config.yaml    # dev defaults
Dockerfile             # single-stage Go build
railway.toml           # healthcheck = /health
```

## Run locally

```sh
make dev          # go run ./cmd/server
make build        # binary at ./bin/server
make test         # go test ./...
make lint         # golangci-lint
make help         # full target list
```

The site listens on `:8080` — http://localhost:8080.

## Configuration

Defaults live in `configs/config.yaml` and `internal/platform/config/config.go`. Override at runtime with `PORTFOLIO_`-prefixed env vars (underscores are key separators):

```sh
PORTFOLIO_SERVER_PORT=9000
PORTFOLIO_EMAIL_PROVIDER=resend
PORTFOLIO_EMAIL_RESEND_API_KEY=re_...
```

Railway's generic `PORT` is also bound directly.

For local development, copy [`.env.example`](./.env.example) to `.env` (gitignored). The server loads it via `godotenv` at startup and Viper picks up the values like any other env override.

### Email providers

The contact form picks a sender based on `email.provider`:

| Value     | Behavior                                          |
|-----------|---------------------------------------------------|
| *(empty)* | Dev mode — message is logged via slog, no network |
| `smtp`    | net/smtp via STARTTLS on port 587                 |
| `resend`  | Resend HTTP API                                   |

## Deploy

One Docker image deploys to Railway. Reproduce the production image locally:

```sh
make docker-build
make docker-run    # binds :8080 and forwards PORTFOLIO_EMAIL_* vars
```

## Conventions

See [`CLAUDE.md`](./CLAUDE.md) for the project's working notes — the BFF rule (no JSON APIs), Datastar conventions (`{contentType: 'form'}`, `_`-prefixed UI signals), and the two-layered styling system (native CSS tokens + Tailwind token bridge).
