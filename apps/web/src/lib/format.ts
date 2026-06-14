/** Country code (e.g. "JP") → flag emoji 🇯🇵. Empty string when unknown. */
export function countryFlag(cc: string | null | undefined): string {
  if (!cc || cc.length !== 2) return "";
  const base = 0x1f1e6; // regional indicator 'A'
  let out = "";
  for (const ch of cc.toUpperCase()) {
    const code = ch.charCodeAt(0);
    if (code < 65 || code > 90) return "";
    out += String.fromCodePoint(base + (code - 65));
  }
  return out;
}
