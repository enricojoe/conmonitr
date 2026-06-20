package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"conmonitr/backend/internal/docker"
)

// WSLogs follows and streams a single container's logs to the client.
func (h *Handler) WSLogs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tail := r.URL.Query().Get("tail")
	if tail == "" {
		tail = "200"
	}

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

	lines := make(chan docker.LogLine, 64)
	go func() {
		_ = h.svc.StreamLogs(ctx, id, tail, lines)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case line := <-lines:
			if err := conn.WriteJSON(line); err != nil {
				return
			}
		}
	}
}
