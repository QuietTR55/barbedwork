package main

import (
	"backend/internal/di"
	"backend/pkg/middleware"
	"net/http"
	"time"
)

func SetupRoutes(mux *http.ServeMux, container *di.Container) {
	healthCheckStack := []middleware.Middleware{
		middleware.RateLimitMiddleware(container.DefaultLimiter, time.Minute),
	}
	mux.Handle("/api/health", middleware.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}),
		healthCheckStack...,
	))

	container.AdminAuthHandler.RegisterRoutes(mux)
	container.AdminDashboardHandler.RegisterRoutes(mux)
	container.UserHandler.RegisterRoutes(mux)
	container.UserAuthHandler.RegisterRoutes(mux)
	container.WorkspaceHandler.RegisterRoutes(mux)
}
