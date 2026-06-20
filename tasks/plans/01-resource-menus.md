# Plan 01 — Image / Volume / Network menu pages (read-only)

**Owner:** Sonnet subagent A
**Scope decision:** Read-only list views only. No remove/prune/inspect actions.

## Goal
Add three new navigable sections — **Images**, **Volumes**, **Networks** — each a
read-only table/grid listing the resources reported by the Docker engine. Add a top
navigation menu so the user can move between Containers (existing dashboard) and the
three new pages.

## Backend (Go, `backend/`)
Follow the existing layered pattern: `internal/docker` (SDK wrapper + DTOs) →
`internal/api` (HTTP handlers) → `internal/api/router.go` (routes).

1. `internal/docker/images.go`
   - DTO `ImageSummary` (JSON-tagged): `id`, `repoTags []string`, `size int64`,
     `created int64`, `containers int64` (use `-1`→0 if unknown).
   - `func (s *Service) ListImages(ctx) ([]ImageSummary, error)` using
     `s.cli.ImageList(ctx, image.ListOptions{All: false})`
     (`github.com/docker/docker/api/types/image`).
2. `internal/docker/volumes.go`
   - DTO `VolumeSummary`: `name`, `driver`, `mountpoint`, `scope`, `createdAt`.
   - `func (s *Service) ListVolumes(ctx) ([]VolumeSummary, error)` using
     `s.cli.VolumeList(ctx, volume.ListOptions{})`
     (`github.com/docker/docker/api/types/volume`); iterate `resp.Volumes`.
3. `internal/docker/networks.go`
   - DTO `NetworkSummary`: `id`, `name`, `driver`, `scope`, `created` (string),
     `internal bool`.
   - `func (s *Service) ListNetworks(ctx) ([]NetworkSummary, error)` using
     `s.cli.NetworkList(ctx, network.ListOptions{})`
     (`github.com/docker/docker/api/types/network`).
   - **Verify exact return type against the installed SDK (docker v28) before coding;
     compile to confirm.**
4. `internal/api/resources.go` — handlers `ListImages`, `ListVolumes`, `ListNetworks`
   mirroring `ListContainers` (use `writeJSON` / `writeError`, `r.Context()`).
5. `internal/api/router.go` — add under `/api`:
   ```go
   r.Get("/images", h.ListImages)
   r.Get("/volumes", h.ListVolumes)
   r.Get("/networks", h.ListNetworks)
   ```

## Frontend (React/TS, `frontend/src/`)
1. `types.ts` — add `ImageSummary`, `VolumeSummary`, `NetworkSummary` mirroring the Go
   DTOs (nullable arrays where Go may emit `null`).
2. `api/client.ts` — add `listImages()`, `listVolumes()`, `listNetworks()` to the `api`
   object using the existing `req<T>` helper.
3. `components/NavMenu.tsx` — a small top nav with links: Containers (`/`),
   Images (`/images`), Volumes (`/volumes`), Networks (`/networks`). Use
   react-router `NavLink` with an active style; match existing Tailwind/zinc theme.
4. New pages, each following the `Dashboard.tsx` data-fetch pattern (load on mount,
   error state, simple table styled with existing Tailwind classes):
   - `pages/Images.tsx` — columns: repo:tag, image id (short), size (use
     `formatBytes` from `utils/format`), created.
   - `pages/Volumes.tsx` — columns: name, driver, mountpoint, scope.
   - `pages/Networks.tsx` — columns: name, driver, scope, internal.
5. `App.tsx` — add routes `/images`, `/volumes`, `/networks`.
6. Render `<NavMenu />` on the Dashboard and the three new pages (a shared header
   region). **Do NOT modify `pages/ContainerView.tsx`** (owned by Plan 02).

## Constraints
- Reuse `formatBytes`/`formatPercent` from `utils/format.ts`; don't duplicate.
- No new npm or Go dependencies (Docker SDK already vendored).
- Match existing code style, comment density, and Tailwind theme.

## Verification (must pass before reporting done)
- `cd backend && go build ./... && go vet ./...`
- `cd frontend && npm run build` (runs `tsc --noEmit` + vite build)
- Report any type mismatches discovered against the SDK.

## Files this agent may touch (no others)
backend/internal/docker/{images,volumes,networks}.go,
backend/internal/api/resources.go, backend/internal/api/router.go,
frontend/src/types.ts, frontend/src/api/client.ts,
frontend/src/components/NavMenu.tsx,
frontend/src/pages/{Images,Volumes,Networks}.tsx,
frontend/src/App.tsx, frontend/src/pages/Dashboard.tsx

## Results
_(fill in on completion)_
