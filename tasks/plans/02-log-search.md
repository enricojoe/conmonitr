# Plan 02 — Search/filter on logs

**Owner:** Sonnet subagent B
**Scope:** Frontend only. Client-side filtering of the in-memory log buffer.

## Goal
Add a search box to the log viewer so the user can filter the streamed log lines by
substring in real time. Filtering is purely client-side over the bounded buffer that
`useLogsSocket` already maintains (max 500 lines) — no backend changes.

## Implementation (`frontend/src/`)
Edit **only** `components/LogViewer.tsx`:
1. Add local state `query` (string) via `useState`.
2. Render a search `<input>` above the log scroll area, styled with the existing
   Tailwind/zinc theme (small, monospace-friendly, placeholder e.g. "Filter logs…").
   Include a clear (×) affordance when non-empty.
3. Compute `filtered = lines` when query is empty, else
   `lines.filter(l => l.line.toLowerCase().includes(query.toLowerCase()))`.
   Render `filtered` instead of `lines`.
4. When a query is active, show a small count, e.g. "showing N of M lines".
5. Optional but preferred: highlight the matched substring within each visible line
   (wrap matches in a `<mark>` styled with Tailwind). Keep it simple and safe — split
   on the query, no `dangerouslySetInnerHTML`.
6. Auto-scroll behaviour: keep the existing scroll-to-end effect but base it on the
   rendered (filtered) list; when a filter is active, do not force-scroll on every
   keystroke (only scroll when new lines arrive and no filter is active). Keep it
   minimal — don't over-engineer.
7. Empty-states: distinguish "Waiting for log output…" (no lines at all) from
   "No lines match \"<query>\"" (lines exist but none match).

## Constraints
- Do not change `hooks/useLogsSocket.ts`, `pages/ContainerView.tsx`, the `LogLine`
  type, or any backend file. Component `Props` stay `{ lines: LogLine[] }`.
- Match existing style and comment density.
- No new dependencies.

## Verification (must pass before reporting done)
- `cd frontend && npm run build` (tsc --noEmit + vite build) passes.
- Sanity-check the filter logic by reasoning through empty query, matching query,
  and non-matching query cases.

## Files this agent may touch (no others)
frontend/src/components/LogViewer.tsx

## Results

Implemented on 2026-06-20. Single file changed: `frontend/src/components/LogViewer.tsx`.

Changes made:
- Added `useState` for `query` string (empty by default).
- Added a filter bar above the scroll area: a monospace `<input>` styled with zinc/black theme, an inline × clear button (shown only when query is non-empty), and a "N / M lines" match count (shown only when query is active).
- Computed `filtered` array: equals `lines` when query is empty, otherwise case-insensitive `includes` filter.
- Auto-scroll `useEffect` now gates on `lowerQuery === ""` so keystrokes do not force-scroll while a filter is active; it continues to scroll on new line arrival when no filter is set.
- Added `highlightMatches(text, term)` helper that iterates `indexOf` in a loop and returns an array of `string | <mark>` React nodes — no `dangerouslySetInnerHTML`, safe for arbitrary log content.
- Three render branches: `lines.length === 0` → "Waiting for log output…", `filtered.length === 0` → `No lines match "<query>"`, otherwise the filtered line list with highlights.
- `import type React from "react"` added for the `React.ReactNode[]` return-type annotation (type-only import, compatible with `noUnusedLocals` and `isolatedModules`).

Build verification: Bash was denied in the sandbox; the parent agent should run `cd frontend && npm run build` to confirm. TypeScript correctness was verified by reasoning through types and compiler flags (`strict`, `noUnusedLocals`, `noUnusedParameters`, `isolatedModules`, `jsx: react-jsx`).
