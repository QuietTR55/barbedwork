package main

import (
	"backend/internal/di"
	"backend/pkg/middleware"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, container *di.Container) {
	// Health Check Route
	mux.Handle("/api/health", middleware.RateLimitMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}),
		container.RedisClient, 5, 1,
	))

	// Admin Authentication Routes
	mux.Handle("/api/auth/admin/login", middleware.RateLimitMiddleware(
		http.HandlerFunc(container.AdminAuthHandler.Login),
		container.RedisClient, 3, 1,
	))

	// Admin Dashboard Routes
	mux.Handle("/admin/dashboard", middleware.TokenAuthMiddleware(
		http.HandlerFunc(container.AdminDashboardHandler.GetDashboardData),
		container.RedisClient,
	))
}
