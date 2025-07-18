package handlers

import (
	"backend/internal/services"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"net/http"
	"time"
)

type UserHandler struct {
	userService *services.UserService
	store       utilities.SessionStore
	limiter     ratelimiter.RateLimiter
}

func NewUserHandler(userService *services.UserService, store utilities.SessionStore, limiter ratelimiter.RateLimiter) *UserHandler {
	return &UserHandler{userService: userService, store: store, limiter: limiter}
}

func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {

	userStack := []middleware.Middleware{
		middleware.TokenAuthMiddleware(h.store),
		middleware.RateLimitMiddleware(h.limiter, time.Minute, "user_profile"),
	}

	router.Handle("/user", middleware.Chain(
		http.HandlerFunc(h.GetUser),
		userStack...,
	))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := utilities.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
}