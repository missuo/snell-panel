# Snell Panel

Manage Snell proxy nodes and generate subscription links. Backend is **Hono on
Cloudflare Workers + D1**; the panel is a **HeroUI v3 SPA** served by the same Worker.
Tooling is **Bun**. See [`docs/DESIGN.md`](docs/DESIGN.md) for the full design.

## Layout

```
apps/server   Hono Worker (API + serves the SPA) + D1 schema/migrations
apps/web      Vite + React + HeroUI v3 SPA  → builds to apps/web/dist
packages/shared  Shared TS types + zod schemas
scripts       snell-install.sh (installer) + import-legacy.ts
```

## Develop

```bash
bun install

# terminal 1 — Worker + local D1
cd apps/server
cp .dev.vars.example .dev.vars      # set ACCESS_TOKEN / API_TOKEN
bun run db:migrate:local
bun run dev                          # http://localhost:8787

# terminal 2 — SPA (proxies /api to the Worker)
cd apps/web
bun run dev                          # http://localhost:5173
```

## Configuration

| Name | Kind | Purpose |
|---|---|---|
| `ACCESS_TOKEN` | secret | Panel login (control plane) |
| `API_TOKEN` | secret | Data-plane master write token (never leaves the backend) |
| `SNELL_V5_VERSION` | var | Exact V5 build (default `v5.0.1`) |
| `SNELL_V6_VERSION` | var | Exact V6 build (default `v6.0.0b2`) |

The **subscription token** is separate, stored in D1, and rotatable from the panel
(Subscription → Reset token). It is independent of `ACCESS_TOKEN`.

## Deploy (Cloudflare)

```bash
cd apps/server
wrangler login                                   # or set CLOUDFLARE_API_TOKEN
wrangler d1 create snell-panel                    # paste database_id into wrangler.jsonc
wrangler secret put ACCESS_TOKEN
wrangler secret put API_TOKEN
wrangler d1 migrations apply snell-panel --remote
bun run --filter '@snell-panel/web' build         # produces apps/web/dist
wrangler deploy
```

## Node lifecycle

1. **Add Node** in the panel — pick V5/V6, name, optionally pre-fill IP/Port → a
   `pending` node.
2. **Install** — copy the generated command, run it on the server. It installs
   snell, then registers `ip/port/psk` back → the node becomes `active`.
3. **Relay** — clone an active node behind a different IP/port (transit front).
4. **Upgrade** — migrate a V4/V5 node to V6 in place (config migration + binary swap +
   re-report).

## Installer (`scripts/snell-install.sh`)

```bash
bash <(curl -fsSL https://panel/install.sh) install \
  --api-url https://panel --node-id <id> --token <one-time> --version 6
bash <(curl -fsSL https://panel/install.sh) upgrade   --token <one-time> ...
bash <(curl -fsSL https://panel/install.sh) uninstall [--api-token <API_TOKEN>]
```

The installer stores `node_id` in `/etc/snell/.install_meta`; `uninstall` deletes the
panel entry **by node id** (not IP).

## Import from the legacy panel

```bash
bun scripts/import-legacy.ts "https://old-panel/entries?token=..." > import.sql
cd apps/server && wrangler d1 execute snell-panel --remote --file=../../import.sql
```

Drops V5 nodes, preserves each `node_id`, and re-assigns integer ids from 1.
