import { Hono } from "hono";
import type { AppEnv } from "../env";
// Bundled as a text module via wrangler `rules` (type: "Text", globs: ["**/*.sh"]).
import installScript from "../../../../scripts/snell-install.sh";

const router = new Hono<AppEnv>();

// GET /install.sh — serve the installer so the generated command can curl it.
router.get("/", (c) => {
  c.header("Content-Type", "text/x-shellscript; charset=utf-8");
  c.header("Cache-Control", "no-cache");
  return c.body(installScript);
});

export default router;
