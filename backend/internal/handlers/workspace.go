package handlers

import (
	"backend/internal/repos"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"net/http"
	"time"
)

type WorkspaceHandler struct {
	workspaceRepo    *repos.WorkspaceRepo
	store            utilities.SessionStore
	limiter          ratelimiter.RateLimiter
}

func NewWorkspaceHandler(workspaceRepo *repos.WorkspaceRepo, store utilities.SessionStore, limiter ratelimiter.RateLimiter) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceRepo: workspaceRepo,
		store:        store,
		limiter:     limiter,
	}
}

func (h *WorkspaceHandler) RegisterRoutes(router *http.ServeMux) {
	stack := []middleware.Middleware{
		middleware.TokenAuthMiddleware(h.store),
		middleware.RateLimitMiddleware(h.limiter, time.Minute),
	}
	router.Handle("/api/workspaces", middleware.Chain(
		http.HandlerFunc(h.GetWorkspaces),
		stack...,
	))
}


func (h *WorkspaceHandler) GetWorkspaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := utilities.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	workspaces, err := h.workspaceRepo.GetUserWorkspaces(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get workspaces", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaces)
}
