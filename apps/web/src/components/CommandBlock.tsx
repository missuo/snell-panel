import { useState } from "react";
import { Highlight, themes } from "prism-react-renderer";

/** A one-line bash command rendered with syntax highlighting. Click to copy. */
export function CommandBlock({ code }: { code: string }) {
  const [copied, setCopied] = useState(false);
  const dark = document.documentElement.classList.contains("dark");

  async function copy() {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      /* clipboard unavailable */
    }
  }

  return (
    <div className="relative">
      <button
        type="button"
        onClick={copy}
        title="Click to copy"
        className="block w-full cursor-pointer rounded-xl bg-surface p-3 text-left ring-1 ring-black/10 transition hover:ring-2 hover:ring-accent dark:ring-white/10"
      >
        <Highlight
          code={code}
          language="bash"
          theme={dark ? themes.vsDark : themes.github}
        >
          {({ tokens, getLineProps, getTokenProps }) => (
            <code className="block overflow-x-auto font-mono text-xs whitespace-pre">
              {tokens.map((line, i) => (
                <span key={i} {...getLineProps({ line })}>
                  {line.map((token, key) => (
                    <span key={key} {...getTokenProps({ token })} />
                  ))}
                </span>
              ))}
            </code>
          )}
        </Highlight>
      </button>
      {copied && (
        <span className="absolute top-2 right-2 rounded bg-success px-2 py-0.5 text-xs text-success-foreground">
          Copied!
        </span>
      )}
    </div>
  );
}
