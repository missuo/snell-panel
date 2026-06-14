import { Hono } from "hono";
import { cors } from "hono/cors";
import { createDb } from "./db/client";
import type { AppEnv } from "./env";
import { requireAccess } from "./middleware/auth";
import { resolveVersions } from "./lib/versions";
import nodesRouter from "./routes/nodes";
import registerRouter from "./routes/register";
import subscribeRouter from "./routes/subscribe";
import settingsRouter from "./routes/settings";
import installRouter from "./routes/install";

const app = new Hono<AppEnv>();

// Attach a per-request Drizzle client.
app.use("*", async (c, next) => {
  c.set("db", createDb(c.env));
  await next();
});

// CORS for the API. Same-origin in production; cross-origin during `vite dev`.
// Auth is bearer/query-token based (no cookies), so allowing any origin is safe.
app.use("/api/*", cors());

app.get("/api/snell-versions", requireAccess, (c) => c.json(resolveVersions(c.env)));

// Installer callback (token / API_TOKEN auth) must be registered before the
// admin router so its per-route auth applies, not the admin guard.
app.route("/api/nodes", registerRouter);
app.route("/api/nodes", nodesRouter);
app.route("/api/subscribe", subscribeRouter);
app.route("/api/settings", settingsRouter);
app.route("/install.sh", installRouter);

// Defensive SPA fallback. With `run_worker_first` scoped to /api/* and
// /install.sh, other paths are served by the asset layer and never reach here.
app.all("*", async (c) => {
  if (c.env.ASSETS) return c.env.ASSETS.fetch(c.req.raw);
  return c.notFound();
});

export default app;
