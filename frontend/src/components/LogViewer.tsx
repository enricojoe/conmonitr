import { useEffect, useRef, useState } from "react";
import type React from "react";
import type { LogLine } from "../types";

interface Props {
  lines: LogLine[];
}

// Splits `text` around all case-insensitive occurrences of `term` and returns
// an array of React nodes with matches wrapped in a highlighted <mark>.
function highlightMatches(text: string, term: string): React.ReactNode[] {
  if (!term) return [text];
  const lower = text.toLowerCase();
  const lowerTerm = term.toLowerCase();
  const nodes: React.ReactNode[] = [];
  let cursor = 0;
  let idx = lower.indexOf(lowerTerm, cursor);
  while (idx !== -1) {
    if (idx > cursor) nodes.push(text.slice(cursor, idx));
    nodes.push(
      <mark key={idx} className="bg-yellow-400/30 text-yellow-200 rounded-sm">
        {text.slice(idx, idx + term.length)}
      </mark>
    );
    cursor = idx + term.length;
    idx = lower.indexOf(lowerTerm, cursor);
  }
  if (cursor < text.length) nodes.push(text.slice(cursor));
  return nodes;
}

// LogViewer renders streaming log lines, supports substring filtering, and
// auto-scrolls to the newest entry when no filter is active.
export default function LogViewer({ lines }: Props) {
  const endRef = useRef<HTMLDivElement>(null);
  const [query, setQuery] = useState("");

  // Compute filtered set; when query is empty the full buffer is shown.
  const lowerQuery = query.toLowerCase();
  const filtered =
    lowerQuery === ""
      ? lines
      : lines.filter((l) => l.line.toLowerCase().includes(lowerQuery));

  // Auto-scroll only when no filter is active so keystrokes don't jump the view.
  useEffect(() => {
    if (lowerQuery === "") {
      endRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [lines, lowerQuery]);

  return (
    <div className="rounded-xl border border-zinc-800 bg-black/60 p-3 space-y-2">
      {/* Search / filter bar */}
      <div className="flex items-center gap-1">
        <div className="relative flex-1">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Filter logs…"
            className="w-full rounded border border-zinc-700 bg-zinc-900 px-2 py-1 font-mono text-xs text-zinc-200 placeholder-zinc-600 outline-none focus:border-zinc-500"
          />
          {query && (
            <button
              onClick={() => setQuery("")}
              aria-label="Clear filter"
              className="absolute right-1.5 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-300"
            >
              ×
            </button>
          )}
        </div>
        {/* Match count shown only when a query is active */}
        {query && (
          <span className="shrink-0 text-xs text-zinc-500">
            {filtered.length} / {lines.length} lines
          </span>
        )}
      </div>

      {/* Log scroll area */}
      <div className="h-80 overflow-auto font-mono text-xs leading-relaxed">
        {lines.length === 0 ? (
          // No lines have arrived from the socket yet.
          <p className="text-zinc-600">Waiting for log output…</p>
        ) : filtered.length === 0 ? (
          // Lines exist but none match the active query.
          <p className="text-zinc-600">
            No lines match &ldquo;{query}&rdquo;
          </p>
        ) : (
          filtered.map((l, i) => (
            <div
              key={i}
              className={l.stream === "stderr" ? "text-rose-300" : "text-zinc-300"}
            >
              <span className="mr-2 text-zinc-600">
                {new Date(l.timestamp).toLocaleTimeString()}
              </span>
              {highlightMatches(l.line, query)}
            </div>
          ))
        )}
        <div ref={endRef} />
      </div>
    </div>
  );
}
