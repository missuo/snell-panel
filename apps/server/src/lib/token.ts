import { and, eq } from "drizzle-orm";
import { installTokens } from "../db/schema";
import type { Db } from "../db/client";

/** One-time install/upgrade tokens live for 5 minutes. */
export const TOKEN_TTL_SECONDS = 5 * 60;

export type TokenPurpose = "install" | "upgrade";

/** 32 random bytes, hex-encoded. */
export function newToken(): string {
  const bytes = new Uint8Array(32);
  crypto.getRandomValues(bytes);
  return Array.from(bytes, (b) => b.toString(16).padStart(2, "0")).join("");
}

export async function mintToken(
  db: Db,
  nodeId: string,
  purpose: TokenPurpose,
  now: number,
): Promise<{ token: string; expiresAt: number }> {
  const token = newToken();
  const expiresAt = now + TOKEN_TTL_SECONDS;
  await db.insert(installTokens).values({ token, nodeId, purpose, expiresAt });
  return { token, expiresAt };
}

/**
 * Validate a one-time token for a node WITHOUT consuming it. Used by the
 * installer's pre-flight check so a doomed install never starts.
 */
export async function validateToken(
  db: Db,
  token: string,
  nodeId: string,
  now: number,
): Promise<{ ok: boolean; reason?: "missing" | "used" | "expired" }> {
  const rows = await db
    .select()
    .from(installTokens)
    .where(and(eq(installTokens.token, token), eq(installTokens.nodeId, nodeId)))
    .limit(1);

  const row = rows[0];
  if (!row) return { ok: false, reason: "missing" };
  if (row.usedAt !== null) return { ok: false, reason: "used" };
  if (row.expiresAt < now) return { ok: false, reason: "expired" };
  return { ok: true };
}

/**
 * Validate a one-time token for a node and consume it (single-use).
 * Returns true only when the token exists, matches the node, is unused, and
 * has not expired.
 */
export async function consumeToken(
  db: Db,
  token: string,
  nodeId: string,
  now: number,
): Promise<boolean> {
  const rows = await db
    .select()
    .from(installTokens)
    .where(and(eq(installTokens.token, token), eq(installTokens.nodeId, nodeId)))
    .limit(1);

  const row = rows[0];
  if (!row) return false;
  if (row.usedAt !== null) return false;
  if (row.expiresAt < now) return false;

  await db
    .update(installTokens)
    .set({ usedAt: now })
    .where(eq(installTokens.token, token));
  return true;
}
