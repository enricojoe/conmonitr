package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"conmonitr/backend/internal/docker"
)

// upgrader upgrades HTTP requests to WebSocket. Origin checks are disabled
// because this is a local-only dashboard served behind the Vite dev proxy.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// readPump drains incoming WS frames purely to detect client disconnect,
// cancelling ctx (via cancel) when the connection closes.
func readPump(conn *websocket.Conn, cancel context.CancelFunc) {
	defer cancel()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

// WSStats streams aggregate live metrics for all running containers. It keeps
// one engine stats stream per running container and emits a JSON array of the
// latest snapshot every two seconds, refreshing the container set periodically.
func (h *Handler) WSStats(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	if !wsAuthenticate(conn) {
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	go readPump(conn, cancel)

	metrics := make(chan docker.Metric, 64)
	done := make(chan string, 16)
	latest := map[string]docker.Metric{}
	active := map[string]context.CancelFunc{}

	startStreamer := func(id, name string) {
		sctx, scancel := context.WithCancel(ctx)
		active[id] = scancel
		go func() {
			_ = h.svc.StreamStats(sctx, id, name, metrics)
			select {
			case done <- id:
			case <-ctx.Done():
			}
		}()
	}

	refresh := func() {
		items, err := h.svc.List(ctx)
		if err != nil {
			return
		}
		for _, c := range items {
			if c.State != "running" {
				continue
			}
			if _, ok := active[c.ID]; !ok {
				startStreamer(c.ID, c.Name)
			}
		}
	}

	refresh()
	refreshTicker := time.NewTicker(5 * time.Second)
	emitTicker := time.NewTicker(2 * time.Second)
	defer refreshTicker.Stop()
	defer emitTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-metrics:
			latest[m.ID] = m
		case id := <-done:
			if cancelFn, ok := active[id]; ok {
				cancelFn()
				delete(active, id)
			}
			delete(latest, id)
		case <-refreshTicker.C:
			refresh()
		case <-emitTicker.C:
			snapshot := make([]docker.Metric, 0, len(latest))
			for _, m := range latest {
				snapshot = append(snapshot, m)
			}
			if err := conn.WriteJSON(snapshot); err != nil {
				return
			}
		}
	}
}

// WSStatsOne streams high-resolution metrics for a single container.
func (h *Handler) WSStatsOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	if !wsAuthenticate(conn) {
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	go readPump(conn, cancel)

	name := id
	if d, err := h.svc.Inspect(ctx, id); err == nil {
		name = d.Name
	}

	metrics := make(chan docker.Metric, 16)
	go func() {
		_ = h.svc.StreamStats(ctx, id, name, metrics)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-metrics:
			if err := conn.WriteJSON(m); err != nil {
				return
			}
		}
	}
}
