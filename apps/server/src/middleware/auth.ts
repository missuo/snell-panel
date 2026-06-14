import type { Context } from "hono";
import { createMiddleware } from "hono/factory";
import type { AppEnv } from "../env";

/** Constant-time string comparison. */
export function safeEqual(a: string, b: string): boolean {
  if (a.length !== b.length) return false;
  let out = 0;
  for (let i = 0; i < a.length; i++) out |= a.charCodeAt(i) ^ b.charCodeAt(i);
  return out === 0;
}

/** Pull a token from `Authorization: Bearer <t>` or the `?token=` query param. */
export function extractToken(c: Context<AppEnv>): string | null {
  const auth = c.req.header("Authorization");
  if (auth) {
    const m = auth.match(/^Bearer\s+(.+)$/i);
    if (m) return m[1].trim();
  }
  return c.req.query("token") ?? null;
}

/** Control-plane guard: requires the panel ACCESS_TOKEN. */
export const requireAccess = createMiddleware<AppEnv>(async (c, next) => {
  const tok = extractToken(c);
  if (!tok || !safeEqual(tok, c.env.ACCESS_TOKEN)) {
    return c.json({ error: "Unauthorized" }, 401);
  }
  await next();
});

/** True when the request carries the master API_TOKEN. */
export function hasApiToken(c: Context<AppEnv>): boolean {
  const tok = extractToken(c);
  return tok !== null && safeEqual(tok, c.env.API_TOKEN);
}

/** Data-plane guard: accepts either the panel ACCESS_TOKEN or the API_TOKEN. */
export const requireAccessOrApiToken = createMiddleware<AppEnv>(async (c, next) => {
  const tok = extractToken(c);
  if (tok && (safeEqual(tok, c.env.ACCESS_TOKEN) || safeEqual(tok, c.env.API_TOKEN))) {
    await next();
    return;
  }
  return c.json({ error: "Unauthorized" }, 401);
});
