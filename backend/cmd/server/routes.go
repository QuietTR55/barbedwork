package main

import (
	"backend/internal/di"
	"backend/pkg/middleware"
	"net/http"
	"time"
)

func SetupRoutes(mux *http.ServeMux, container *di.Container) {
	// Health Check Route
	mux.Handle("/api/health", middleware.RateLimitMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}),
		container.Limiter,
		time.Minute,
	))

	container.AdminAuthHandler.RegisterRoutes(mux)
	container.AdminDashboardHandler.RegisterRoutes(mux)
	container.UserHandler.RegisterRoutes(mux)
}
