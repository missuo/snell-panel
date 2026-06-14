<div align="center">
  <img src="apps/web/public/favicon.svg" width="76" alt="Snell Panel" />
  <h1>Snell Panel</h1>
  <p>Manage Snell proxy nodes and generate subscription links.<br/>
  Hono on <b>Cloudflare Workers + D1</b>, with a <b>HeroUI v3</b> panel served from the same Worker.</p>

  <a href="https://deploy.workers.cloudflare.com/?url=https://github.com/missuo/snell-panel">
    <img src="https://deploy.workers.cloudflare.com/button" alt="Deploy to Cloudflare" />
  </a>
</div>

---

## Features

- **V5 / V6 nodes** with an *add-then-install* flow — create a node in the panel, run the generated one-line command on your server, and it back-fills `ip/port/psk`.
- **Two independent secrets**: a panel **Access Token** and a backend **API Token**; servers only ever receive a **per-node one-time install token**.
- **Relay / transit nodes** — clone an active node behind a different IP/port.
- **Enable / disable** — hide a node from subscriptions while it keeps running.
- **Subscriptions** in **Surge / Shadowrocket / Mihomo**, with a **rotatable** subscribe token and flag / filter / relay options.
- **In-place V4/V5 → V6 upgrade** — validates the PSK, strips removed config keys, swaps the binary, and re-reports.
- **Responsive panel** — dense table on desktop, cards on mobile, light/dark themes.

See [`docs/DESIGN.md`](docs/DESIGN.md) for the full design.

---

## Stack

| Layer | Tech |
|---|---|
| Backend | Hono on Cloudflare Workers, D1 (SQLite) via Drizzle |
| Frontend | Vite + React + HeroUI v3 (Tailwind v4), served as Worker assets |
| Tooling | Bun (workspace, scripts), Wrangler (deploy) |

```
apps/server   Hono Worker (API + serves the SPA) + D1 schema/migrations
apps/web      Vite + React + HeroUI v3 SPA  → builds to apps/web/dist
packages/shared  Shared TS types + zod schemas
scripts       snell-install.sh (installer) + import-legacy.ts
```

---

## One-click deploy

Click the **Deploy to Cloudflare** button above. It forks the repo, provisions the
Worker + D1, and connects [Workers Builds](https://developers.cloudflare.com/workers/ci-cd/builds/)
to your default branch. During setup, use these build settings:

| Setting | Value |
|---|---|
| Build command | `bun install && bun run build` |
| Deploy command | `bunx wrangler deploy --config apps/server/wrangler.jsonc` |

Then finish with the two required post-deploy steps (the button can't know your secrets
or run migrations):

```bash
# 1) set the two panel secrets
printf '%s' "$ACCESS_TOKEN" | bunx wrangler secret put ACCESS_TOKEN --config apps/server/wrangler.jsonc
printf '%s' "$API_TOKEN"    | bunx wrangler secret put API_TOKEN    --config apps/server/wrangler.jsonc

# 2) create the database tables
bunx wrangler d1 migrations apply snell-panel --remote --config apps/server/wrangler.jsonc
```

Generate strong tokens with: `openssl rand -base64 24 | tr -dc 'A-Za-z0-9' | head -c 32; echo`

---

## Manual deploy (CLI)

```bash
git clone https://github.com/missuo/snell-panel && cd snell-panel
bun install

cd apps/server
bunx wrangler login

# create D1, then paste the printed database_id into apps/server/wrangler.jsonc
bunx wrangler d1 create snell-panel
bunx wrangler d1 migrations apply snell-panel --remote

# set the two secrets (independent of each other)
printf '%s' "<access-token>" | bunx wrangler secret put ACCESS_TOKEN
printf '%s' "<api-token>"    | bunx wrangler secret put API_TOKEN

# build the SPA, then deploy the Worker (serves the SPA + API)
cd ../.. && bun run build
cd apps/server && bunx wrangler deploy
```

Open the deployed URL and log in with your **Access Token**.

---

## Local development

```bash
bun install

# terminal 1 — Worker + local D1
cd apps/server
cp .dev.vars.example .dev.vars          # set ACCESS_TOKEN / API_TOKEN
bunx wrangler d1 migrations apply snell-panel --local
bun run dev                             # http://localhost:8787

# terminal 2 — SPA (proxies /api to the Worker)
cd apps/web && bun run dev              # http://localhost:5173
```

---

## Configuration

| Name | Kind | Purpose |
|---|---|---|
| `ACCESS_TOKEN` | secret | Panel login (control plane) |
| `API_TOKEN` | secret | Data-plane master write token (never leaves the backend) |
| `SNELL_V5_VERSION` | var | Exact V5 build (default `v5.0.1`) |
| `SNELL_V6_VERSION` | var | Exact V6 build (default `v6.0.0b2`) |

The **subscription token** is separate, stored in D1, and rotatable from the panel
(Subscription → **Reset token**) — independent of `ACCESS_TOKEN`.

---

## Node lifecycle

1. **Add Node** — pick V5/V6, name, optionally pre-fill IP/Port → a `pending` node.
2. **Install** — copy the generated one-line command, run it on the server; it installs
   snell and registers `ip/port/psk` → the node becomes `active`.
3. **Relay** — clone an active node behind a different IP/port (transit front).
4. **Upgrade** — migrate a V4/V5 node to V6 in place (config migration + binary swap).

The installer (`scripts/snell-install.sh`) stores `node_id` in `/etc/snell/.install_meta`,
so `uninstall` removes the panel entry **by node id**, not by IP.

---

## Import from the legacy panel

```bash
bun scripts/import-legacy.ts "https://old-panel/entries?token=..." > import.sql
cd apps/server && bunx wrangler d1 execute snell-panel --remote --file=../../import.sql
```

Drops V5 nodes, preserves each `node_id`, and re-assigns integer ids from 1.

---

## License

See [`LICENSE`](LICENSE).
