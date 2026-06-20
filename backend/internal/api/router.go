package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter wires all REST and WebSocket routes with logging and CORS.
func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", h.Health)
		r.Route("/containers", func(r chi.Router) {
			r.Get("/", h.ListContainers)
			r.Get("/{id}", h.InspectContainer)
			r.Post("/{id}/start", h.StartContainer)
			r.Post("/{id}/stop", h.StopContainer)
			r.Post("/{id}/restart", h.RestartContainer)
			r.Delete("/{id}", h.RemoveContainer)
		})
		r.Get("/images", h.ListImages)
		r.Get("/volumes", h.ListVolumes)
		r.Get("/networks", h.ListNetworks)
	})

	r.Route("/ws", func(r chi.Router) {
		r.Get("/stats", h.WSStats)
		r.Get("/stats/{id}", h.WSStatsOne)
		r.Get("/logs/{id}", h.WSLogs)
	})

	return r
}
