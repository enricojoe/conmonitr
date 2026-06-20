package api

import (
	"net/http"
)

// ListImages returns all non-intermediate Docker images.
func (h *Handler) ListImages(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListImages(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// ListVolumes returns all Docker volumes.
func (h *Handler) ListVolumes(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListVolumes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// ListNetworks returns all Docker networks.
func (h *Handler) ListNetworks(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListNetworks(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}
