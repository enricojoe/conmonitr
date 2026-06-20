# Plan 03 — .env config for backend and frontend (server/proxy wiring)

**Owner:** Sonnet subagent C
**Scope decision:** Externalize server/proxy wiring only:
- Backend: listen address + optional `DOCKER_HOST`.
- Frontend: backend target host/port used by the Vite dev proxy.

## Goal
Introduce `.env` files in both `backend/` and `frontend/`, load them, and wire the
existing hardcoded values to read from them. Commit `.env.example` templates; keep real
`.env` files git-ignored.

## Backend (`backend/`)
Currently `main.go` reads `CONMONITR_ADDR` via `os.Getenv` and defaults to `:8081`;
the Docker client already honours `DOCKER_HOST` via `client.FromEnv`.

1. Add a **dependency-free** `.env` loader (do NOT add godotenv / any new module — keep
   `go.mod`/`go.sum` unchanged). Create `internal/config/env.go` with
   `func Load(path string)` that:
   - reads the file if it exists (no error if absent),
   - skips blank lines and `#` comments, supports `KEY=VALUE` (trim spaces, strip
     optional surrounding quotes),
   - calls `os.Setenv(key, val)` only if the key is **not already set** in the
     environment (real env vars win over the file).
2. In `main.go`, call `config.Load(".env")` at the very start of `main()` (before
   reading `CONMONITR_ADDR`). Keep the existing `:8081` default. `DOCKER_HOST`, if set
   in `.env`, is consumed by the existing `client.FromEnv`.
3. Create `backend/.env.example`:
   ```
   # Address the backend HTTP/WS server listens on
   CONMONITR_ADDR=:8081
   # Optional: override the Docker engine socket/host
   # DOCKER_HOST=unix:///var/run/docker.sock
   ```
4. Create `backend/.env` with the same default content (git-ignored, for local dev).

## Frontend (`frontend/`)
Vite natively loads `.env` files; only `VITE_`-prefixed vars are exposed to client code,
but `vite.config.ts` can read any var via `loadEnv`.

1. `vite.config.ts` — switch to the function form
   `defineConfig(({ mode }) => { const env = loadEnv(mode, process.cwd(), ""); ... })`
   and derive:
   - dev server `port` from `FRONTEND_PORT` (default `5174`),
   - proxy `target` for `/api` and `/ws` from `BACKEND_URL`
     (default `http://localhost:8081`).
   Import `loadEnv` from `vite`.
2. Create `frontend/.env.example`:
   ```
   # Port the Vite dev server listens on
   FRONTEND_PORT=5174
   # Backend the dev server proxies /api and /ws to
   BACKEND_URL=http://localhost:8081
   ```
3. Create `frontend/.env` with the same default content (git-ignored).

## Housekeeping
1. `.gitignore` — ensure `.env` is ignored but `.env.example` is tracked. Add:
   ```
   # Env files
   .env
   backend/.env
   frontend/.env
   !*.env.example
   ```
   (Adjust to existing structure; do not remove existing entries.)
2. `README.md` — add a short note under "Running" that copying `.env.example` →
   `.env` in each folder configures ports/proxy/Docker host. Keep it brief.

## Constraints
- No new Go modules; `go.mod`/`go.sum` must be unchanged.
- Real environment variables must take precedence over `.env` file values.
- Do not change ports' default behaviour (`:8081` backend, `5174` frontend) when no
  `.env` is present.
- Do NOT edit `main.go` beyond adding the single `config.Load(".env")` call and its
  import. Do not touch router/handler/docker files (owned by Plan 01) or any
  `frontend/src/**` file (owned by Plans 01/02).

## Verification (must pass before reporting done)
- `cd backend && go build ./... && go vet ./...`
- `cd frontend && npm run build` passes.
- Confirm `go.mod`/`go.sum` are unmodified (`git`-less: diff by inspection).
- Reason through precedence: a real `CONMONITR_ADDR` env overrides the `.env` value.

## Files this agent may touch (no others)
backend/internal/config/env.go, backend/main.go,
backend/.env, backend/.env.example,
frontend/vite.config.ts, frontend/.env, frontend/.env.example,
.gitignore, README.md

## Results

Completed 2026-06-20.

### Files created
- `backend/internal/config/env.go` — dependency-free `.env` loader; skips missing files, ignores blank/comment lines, strips optional quotes, uses `os.LookupEnv` guard so real env vars always win.
- `backend/.env` — local dev defaults (`CONMONITR_ADDR=:8081`, `DOCKER_HOST` commented out). Git-ignored.
- `backend/.env.example` — same content, committed as template.
- `frontend/.env` — local dev defaults (`FRONTEND_PORT=5174`, `BACKEND_URL=http://localhost:8081`). Git-ignored.
- `frontend/.env.example` — same content, committed as template.

### Files modified
- `backend/main.go` — added `"conmonitr/backend/internal/config"` import and `config.Load(".env")` as the first line of `main()`. No other changes.
- `frontend/vite.config.ts` — switched to `defineConfig(({ mode }) => …)` function form; imports `loadEnv` from `vite`; derives `port` from `FRONTEND_PORT` (default 5174) and `backendUrl` from `BACKEND_URL` (default `http://localhost:8081`). No new npm deps.
- `.gitignore` — added `.env`, `backend/.env`, `frontend/.env` (ignored) and `!*.env.example` (negation to keep examples tracked).
- `README.md` — added "Configuration" section under "Running" describing both `.env` files and their variables. Also corrected stale port references in the existing run commands (`:8080`→`:8081`, `5173`→`5174`).

### Verification
- Bash was denied in the sandbox; parent must run `go build ./... && go vet ./...` in `backend/` and `npm run build` in `frontend/`.
- `go.mod` and `go.sum` are unchanged — no external module was added; `config` is a new internal package within the existing `conmonitr/backend` module.
- Precedence confirmed by code inspection: `os.LookupEnv` guard in `env.go` ensures real env vars win over file values.

### Parent verification (2026-06-20)
- `backend`: `go build ./...` + `go vet ./...` pass. `go.mod`/`go.sum` unchanged (mtime Jun 14).
- `frontend`: initial `npm run build` failed — `vite.config.ts` used `process.cwd()` but `@types/node` is not installed (`tsc` error TS2580). Fixed by passing `"."` as the `loadEnv` env dir instead of `process.cwd()` (dependency-free). `npm run build` now passes.
