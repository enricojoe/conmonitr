# ConMonitr

A real-time web dashboard that monitors local containerized environments. It interfaces
directly with the Docker engine's Unix socket and streams live performance metrics to a
React frontend over WebSockets.

```
┌──────────────┐   WebSocket (stats, logs)   ┌──────────────┐   Unix socket   ┌──────────────┐
│   React UI   │ ◄─────────────────────────► │  Go backend  │ ◄─────────────► │ Docker engine│
│  (Vite/5173) │   REST (list, lifecycle)    │  (chi/8080)  │   Docker SDK    │              │
└──────────────┘                             └──────────────┘                 └──────────────┘
```

## Features

- **Live metrics** — per-container CPU %, memory, network I/O, and block I/O streamed in real time.
- **Container list + status** — running/stopped state, image, ports, uptime.
- **Lifecycle controls** — start / stop / restart / remove from the UI.
- **Live log streaming** — follow stdout/stderr (correctly demultiplexed via `stdcopy`).
- **Container detail** — networks, mounts/volumes, ports, environment, command.

## Stack

- **Backend:** Go, official Docker SDK (`github.com/docker/docker/client`), `gorilla/websocket`, `go-chi`.
- **Frontend:** Vite + React + TypeScript + Tailwind CSS + Recharts.

## Requirements

- Go 1.22+
- Node 18+
- A running Docker engine (Docker Desktop, OrbStack, etc.) reachable at the default socket
  (`/var/run/docker.sock`) or via `DOCKER_HOST`.

## Running

### Configuration

Both the backend and frontend load a `.env` file from their respective directories
on startup.  Copy the provided example and edit as needed:

```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

Key variables:

| File | Variable | Default | Purpose |
|------|----------|---------|---------|
| `backend/.env` | `CONMONITR_ADDR` | `:8081` | Address the backend HTTP/WS server listens on |
| `backend/.env` | `DOCKER_HOST` | _(system default)_ | Override the Docker engine socket/host |
| `frontend/.env` | `FRONTEND_PORT` | `5174` | Port the Vite dev server listens on |
| `frontend/.env` | `BACKEND_URL` | `http://localhost:8081` | Backend origin the dev proxy forwards to |

Real environment variables always take precedence over `.env` file values.
The `.env` files are git-ignored; `.env.example` files are committed.

### Backend

```bash
cd backend
go run .
# listens on :8081 (override with CONMONITR_ADDR)
```

### Frontend

```bash
cd frontend
npm install
npm run dev
# open http://localhost:5174 (proxies /api and /ws to :8081)
```

Or start both at once from the repo root:

```bash
./dev.sh
```

## API

| Method | Path                             | Description                          |
| ------ | -------------------------------- | ------------------------------------ |
| GET    | `/api/health`                    | Liveness + Docker API version        |
| GET    | `/api/containers`                | List all containers                  |
| GET    | `/api/containers/{id}`           | Inspect a container                  |
| POST   | `/api/containers/{id}/start`     | Start                                |
| POST   | `/api/containers/{id}/stop`      | Stop                                 |
| POST   | `/api/containers/{id}/restart`   | Restart                              |
| DELETE | `/api/containers/{id}?force=`    | Remove                               |
| WS     | `/ws/stats`                      | Aggregate live metrics (all running) |
| WS     | `/ws/stats/{id}`                 | Single-container metrics             |
| WS     | `/ws/logs/{id}?tail=`            | Follow container logs                |

## Notes

- This is a **local, single-host** dashboard with **no authentication** — it binds to localhost
  and exposes container controls. Do not expose it to untrusted networks.
- CPU % is computed from raw engine stat deltas using the same formula as `docker stats`.
