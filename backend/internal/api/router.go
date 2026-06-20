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
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	// Public
	limiter := newLoginLimiter()
	r.Post("/api/auth/login", limiter.middleware(h.Login))
	r.Get("/api/health", h.Health)

	// Protected REST routes — require Authorization: Bearer header
	r.Group(func(r chi.Router) {
		r.Use(JWTAuth)

		r.Route("/api/containers", func(r chi.Router) {
			r.Get("/", h.ListContainers)
			r.Get("/{id}", h.InspectContainer)
			r.Post("/{id}/start", h.StartContainer)
			r.Post("/{id}/stop", h.StopContainer)
			r.Post("/{id}/restart", h.RestartContainer)
			r.Delete("/{id}", h.RemoveContainer)
		})
		r.Get("/api/images", h.ListImages)
		r.Get("/api/volumes", h.ListVolumes)
		r.Get("/api/networks", h.ListNetworks)
	})

	// WebSocket routes — auth via first message (token in URL would appear in logs)
	r.Get("/ws/stats", h.WSStats)
	r.Get("/ws/stats/{id}", h.WSStatsOne)
	r.Get("/ws/logs/{id}", h.WSLogs)

	return r
}
