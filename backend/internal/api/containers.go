package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"conmonitr/backend/internal/docker"
)

// Handler holds shared dependencies for HTTP/WS handlers.
type Handler struct {
	svc *docker.Service
}

// NewHandler creates an API handler backed by the given Docker service.
func NewHandler(svc *docker.Service) *Handler {
	return &Handler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// Health reports liveness and the negotiated Docker API version.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	version, err := h.svc.Ping(r.Context())
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "docker unreachable: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":        "ok",
		"dockerVersion": version,
	})
}

// ListContainers returns all containers.
func (h *Handler) ListContainers(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// InspectContainer returns detailed info for one container.
func (h *Handler) InspectContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	detail, err := h.svc.Inspect(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

// StartContainer starts a stopped container.
func (h *Handler) StartContainer(w http.ResponseWriter, r *http.Request) {
	h.lifecycle(w, r, "started", h.svc.Start)
}

// StopContainer stops a running container.
func (h *Handler) StopContainer(w http.ResponseWriter, r *http.Request) {
	h.lifecycle(w, r, "stopped", h.svc.Stop)
}

// RestartContainer restarts a container.
func (h *Handler) RestartContainer(w http.ResponseWriter, r *http.Request) {
	h.lifecycle(w, r, "restarted", h.svc.Restart)
}

// lifecycle runs a start/stop/restart action and writes a uniform response.
func (h *Handler) lifecycle(w http.ResponseWriter, r *http.Request, verb string, action func(context.Context, string) error) {
	id := chi.URLParam(r, "id")
	if err := action(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": verb, "id": id})
}

// RemoveContainer deletes a container (force via ?force=true).
func (h *Handler) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	force := r.URL.Query().Get("force") == "true"
	if err := h.svc.Remove(r.Context(), id, force); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed", "id": id})
}
