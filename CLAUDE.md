# CLAUDE.md

Working notes for Claude Code in this repo. Read first; follow strictly.

## What this is

Personal portfolio site. Single Go binary, server-rendered HTML, no SPA, no
database, no worker, no public API. Embedded templates and static assets
ship inside the binary; deployment is one Docker image to Railway.

## Stack

- **Go 1.22** — `cmd/server/main.go` (Cobra `serve` subcommand pattern)
- **chi router** + standard chi middleware (`RequestID`, `RealIP`, `Recoverer`, `Compress`) plus our own `StructuredLogger` and `SecurityHeaders`
- **`log/slog`** JSON output to stdout
- **html/template** — parsed once, shared between handlers via `handlers.ParseTemplates`
- **Tailwind v4** browser CDN (`@tailwindcss/browser@4`) — see styling rules below
- **Datastar** (`cdn.jsdelivr.net/gh/starfederation/datastar@main`) — all client interactivity
- **Viper + Cobra** for config and CLI
- **Resend / SMTP / log** mailer providers behind a single `Sender` interface

## Layout

```
cmd/server/main.go               # cobra root + `serve` subcommand
internal/
  platform/
    config/    config.go         # viper, PORTFOLIO_ env prefix, Railway PORT bind
    server/    server.go         # chi router build + graceful shutdown
    middleware/middleware.go     # slog logger + security headers
    notify/    notify.go         # Mailer + Sender interface
               smtp.go           #   - net/smtp via STARTTLS
               resend.go         #   - Resend HTTP API
               log.go            #   - dev-mode slog
    sse/       sse.go            # hand-rolled Datastar SSE writer
  handlers/
    page.go                      # GET / index renderer
    contact.go                   # POST /contact (BFF, SSE response)
    health.go                    # GET /health for Railway probe
    templates.go                 # ParseTemplates() — shared template set
  content/projects.go            # static project list (no repository abstraction)
web/
  embed.go                       # //go:embed templates/*.html static
  templates/
    index.html                   # main page
    contact.html                 # named fragments: contact-form-area, contact-thanks
  static/css/styles.css          # native CSS, design tokens, all components
configs/config.yaml              # dev defaults; secrets come from env
Makefile                         # `make help` for the full target list
Dockerfile                       # single-stage Go build
railway.toml                     # healthcheckPath = "/health"
```

When adding code, **flat platform layout** — no DDD layers (`domain/`, `application/`,
`infrastructure/`, `interfaces/`). One bounded context, no abstraction theatre.

## BFF rule (non-negotiable)

Every endpoint serves only this frontend. **No JSON APIs.** No "we might
expose this externally later." The contract is the rendered HTML + SSE
events Datastar consumes.

- Form submissions: accept `application/x-www-form-urlencoded` (Datastar's `@post(url, {contentType: 'form'})`)
- Responses: stream `text/event-stream` events via `internal/platform/sse`
- Reject requests missing the `Datastar-Request: true` header — that header is the canonical proof a request came from Datastar; it stops anyone from accidentally building a public API on top of these endpoints
- Re-render the relevant template fragment (defined in `web/templates/contact.html`-style files) and patch it into the page using `sse.New(w).PatchElements(...)`. Default mode is `outer`; the target is found by matching the fragment's outer `id` attribute
- On validation failure, re-render the same fragment with the user's values preserved + an inline error — never a separate error page or JSON `{error}`

Pattern reference: `internal/handlers/contact.go` is the canonical example.

## Datastar conventions

Use Datastar for **all** client-side state and behavior. No jQuery, no
inline ad-hoc `<script>` blocks for state, no React-style component
libraries. The two exceptions:

1. The pre-paint theme bootstrap script (avoids FOUC before Datastar loads)
2. The smooth-scroll `<script>` at the bottom of `index.html` (no signal involved — purely a one-shot DOM side effect, not worth Datastar verbosity)

Everything else uses Datastar attributes:

- **Signals**: `data-signals="{themeMode: 'light'}"` on a parent element
- **Persistence**: read/write `localStorage` via `data-effect`
- **Bind to attribute**: `data-attr-data-theme="$themeMode"`, `data-attr-aria-pressed="$themeMode === 'dark'"`
- **Click handler**: `data-on-click="$themeMode = $themeMode === 'dark' ? 'light' : 'dark'"`
- **Form submit (BFF)**: `data-on-submit__prevent="@post('/contact', {contentType: 'form'})"`
- **Server response**: SSE `datastar-patch-elements` events (use `internal/platform/sse`)

Datastar uses `MutationObserver` so server-injected fragments are auto-wired
without re-bootstrapping.

## Design system + Tailwind rules

The styling is intentionally **two-layered**:

### Layer 1 — `web/static/css/styles.css` (native CSS, no Tailwind directives)

Loaded synchronously via `<link rel="stylesheet">`. This is the source of
truth for tokens and components. Every value goes through a CSS variable
defined at `:root` / `[data-theme="dark"]`. **No raw hex, rem, or px
literals in component blocks** — always `var(--token)`.

Tokens defined here (light + dark variants):
- Colors: `--bg`, `--bg-band`, `--bg-band-deep`, `--bg-elev`, `--fg`, `--fg-muted`, `--accent`, `--accent-fg`, `--rule`, `--silhouette`, `--hero-sky`, `--hero-cloud`, `--hero-sun`, `--hero-text`
- Raw palette (rarely use directly): `--sky`, `--sky-deep`, `--cream`, `--paper`, `--paper-soft`, `--ink`, `--ink-soft`, `--ink-deep`, `--cloud`
- Fonts: `--font-display`, `--font-serif`, `--font-sans`, `--font-mono`
- Type scale: `--text-eyebrow`, `--text-h-md`, `--text-h-2xl`, `--text-section`, `--text-cta`
- Tracking: `--tracking-mono` (0.04em), `--tracking-eyebrow` (0.16em)
- Spacing: `--spacing-section`, `--spacing-section-sm`, `--spacing-page-x`
- Layout: `--container-page`
- Radius: `--radius-pill`
- Easing: `--ease-editorial`

### Layer 2 — Inline `<style type="text/tailwindcss">` in `index.html`

The Tailwind v4 browser CDN **only processes `<style type="text/tailwindcss">` blocks** — it ignores `<link rel="stylesheet">` regardless of `type=`. This block uses `@theme inline { --color-bg: var(--bg); ... }` to mirror the same tokens so utilities (`bg-bg`, `text-fg-muted`, `font-display`, `text-h-2xl`, `tracking-eyebrow`, `animate-float`, etc.) resolve to design-system values. Keep this block small — it's a token bridge, not the styling itself.

### When to use what

- **Layout/spacing that affects above-the-fold paint** → CSS component class. Tailwind utilities don't apply until the CDN script executes (~50–200ms), so utility-only spacing causes layout shift.
- **One-off inline tweaks below the fold or post-load** → Tailwind utility (`pt-2`, `mb-7`, `text-fg-muted`, etc.) is fine.
- **Dark mode overrides** → `[data-theme="dark"] .foo { ... }` in CSS, or `dark:foo` Tailwind variant. The `@custom-variant dark` is wired in the inline block.
- **A new compound component (nav, hero, project card, …)** → add as a class in `styles.css`, all values via tokens. Don't use `@apply` (the file is plain CSS, not Tailwind-processed).

### Adding a new design token

1. Add the CSS variable to `:root` (and `[data-theme="dark"]` if it flips) in `styles.css`
2. Mirror it in the inline `<style type="text/tailwindcss">` `@theme inline` block in `index.html` if it should be a Tailwind utility
3. Use it via `var(--name)` in CSS or via the generated utility (`bg-name`, `text-name`, `font-name`, etc.) in templates

## Adding a new page

1. New handler in `internal/handlers/foo.go` — accept the shared `*template.Template`
2. Wire route in `internal/platform/server/server.go` + add to `server.Deps`
3. Build deps in `cmd/server/main.go`
4. New template in `web/templates/foo.html` — uses the same head fragment patterns from `index.html`
5. Render via `tmpl.ExecuteTemplate(w, "foo.html", data)`

## Adding a new BFF interaction

1. Add a Datastar attribute to the relevant element (`data-on-click="@post('/foo', {contentType: 'form'})"`)
2. Define a named fragment in a `web/templates/*.html` file with `{{define "foo-fragment"}}<div id="foo-area">…</div>{{end}}`
3. Include it from the parent template: `{{template "foo-fragment" .Foo}}`
4. Handler: validate → call service → render fragment to a `bytes.Buffer` → `sse.New(w).PatchElements("", sse.ModeOuter, buf.String())`
5. Handler MUST guard with `if r.Header.Get("Datastar-Request") == "" { http.Error(...); return }`

## Mailer

`internal/platform/notify` switches on `cfg.Email.Provider`:

- `""` → `log.go` (dev: prints via slog, no network)
- `"smtp"` → `smtp.go` (net/smtp STARTTLS port 587)
- `"resend"` → `resend.go` (Resend HTTP API)

In Railway: set `PORTFOLIO_EMAIL_PROVIDER=resend` and `PORTFOLIO_EMAIL_RESEND_API_KEY=re_…`. Never commit secrets to `configs/config.yaml`.

## Config

Viper with env prefix `PORTFOLIO_`. Underscores are key separators
(`PORTFOLIO_SERVER_PORT`, `PORTFOLIO_EMAIL_RESEND_API_KEY`). Railway's
generic `PORT` is also bound directly. Defaults live in
`internal/platform/config/config.go::setDefaults`.

## Commands

`make help` lists everything. Most-used:
- `make dev`           — `go run ./cmd/server`
- `make build`         — `go build -o bin/server ./cmd/server`
- `make test` / `make lint` / `make fmt` / `make tidy`
- `make docker-build` / `make docker-run` — reproduces the Railway image locally

## Things to avoid

- Don't add npm / Node / Vite / a JS framework. We removed that intentionally.
- Don't reintroduce DDD layers (`domain/`, `application/`, `infrastructure/`).
- Don't add JSON API endpoints. If something needs server data, render it into the page or stream it as a Datastar SSE patch.
- Don't write hex colors, raw `px` margins, or font families inline in templates or component CSS — go through `var(--token)` or a Tailwind utility backed by the token bridge.
- Don't use `@apply` in `styles.css` — the file is loaded as plain CSS, not Tailwind-processed.
- Don't bypass the `Datastar-Request` header guard on BFF endpoints.
- Don't add a `<script>` tag for state/event-handling — use Datastar attributes.
