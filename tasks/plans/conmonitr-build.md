# ConMonitr Build — Detailed Checklist

Real-time container monitoring dashboard. Go backend (Docker SDK + WebSocket) + React/Vite/Tailwind/Recharts frontend.

## Backend
- [ ] `go mod init` + dependencies (docker SDK, gorilla/websocket, chi, chi/cors)
- [ ] `internal/docker/client.go` — Docker client wrapper
- [ ] `internal/docker/containers.go` — List, Inspect, Start, Stop, Restart, Remove
- [ ] `internal/docker/stats.go` — stats stream + CPU% computation
- [ ] `internal/docker/logs.go` — log follow + stdcopy demux
- [ ] `internal/api/router.go` — chi routes + CORS
- [ ] `internal/api/containers.go` — REST handlers
- [ ] `internal/api/ws_stats.go` — WS aggregate + per-container stats
- [ ] `internal/api/ws_logs.go` — WS logs
- [ ] `main.go` — wiring + graceful shutdown
- [ ] `go build ./...` + `go vet ./...` clean
- [ ] curl + WS smoke tests pass

## Frontend
- [ ] Vite React-TS scaffold
- [ ] Tailwind + PostCSS + Recharts + react-router installed
- [ ] vite.config proxy for /api and /ws
- [ ] types.ts + api/client.ts
- [ ] hooks: useStatsSocket, useLogsSocket
- [ ] components: ContainerGrid, StatusBadge, ActionButtons, MetricChart, ContainerDetail, LogViewer
- [ ] pages: Dashboard, ContainerView + App routing
- [ ] `npm run build` clean
- [ ] manual e2e: list, sparklines, lifecycle, detail charts, logs

## Docs
- [ ] root README.md with run instructions
- [ ] root dev script / Makefile (optional)

## Results
Completed 2026-06-14. Backend builds clean (`go build`/`go vet`) and all endpoints verified
end-to-end against the live OrbStack engine (Docker API 1.54): health, list, inspect, restart,
remove, aggregate stats WS (CPU%/mem/net/blk computed correctly), per-container stats WS, and
log streaming (stdcopy demux confirmed — nginx logs arrive tagged `stderr`). Frontend type-checks
and builds (`npm run build`). Full-stack integration verified through the Vite dev proxy: REST and
WebSocket (handshake 101 + live frames) both proxy correctly from :5173 → :8080. README + dev.sh
added. No deviations from plan; chose Tailwind v4 (zero-config) and pinned React 18 for Recharts
peer compatibility.
