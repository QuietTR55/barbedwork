package handlers

import (
	"backend/internal/repos"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WorkspaceHandler struct {
	workspaceRepo    *repos.WorkspaceRepo
	store            utilities.SessionStore
	limiter          ratelimiter.RateLimiter
	permissionChecker *utilities.PermissionChecker
}

func NewWorkspaceHandler(workspaceRepo *repos.WorkspaceRepo, store utilities.SessionStore, limiter ratelimiter.RateLimiter, permissionChecker *utilities.PermissionChecker) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceRepo:   workspaceRepo,
		store:           store,
		limiter:        limiter,
		permissionChecker: permissionChecker,
	}
}

func (h *WorkspaceHandler) RegisterRoutes(router *http.ServeMux) {
	stack := []middleware.Middleware{
		middleware.TokenAuthMiddleware(h.store),
		middleware.RateLimitMiddleware(h.limiter, time.Minute, "workspace_read"),
	}

	modifcationStack:= []middleware.Middleware{
		middleware.TokenAuthMiddleware(h.store),
		middleware.RateLimitMiddleware(h.limiter, time.Minute, "workspace_modify"),
		middleware.PermissionMiddleware(h.permissionChecker, "admin", "moderator"),
	}

	router.Handle("/api/workspaces", middleware.Chain(
		http.HandlerFunc(h.GetWorkspaces),
		stack...,
	))
	router.Handle("/api/workspaces/{workspaceId}", middleware.Chain(
		http.HandlerFunc(h.GetWorkspace),
		stack...,
	))
	router.Handle("/api/workspaces/{workspaceId}/channels", middleware.Chain(
		http.HandlerFunc(h.CreateChannel),
		modifcationStack...,
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

func (h *WorkspaceHandler) GetWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	workspaceID := r.PathValue("workspaceId")
	fmt.Println("Workspace ID:", workspaceID)
	if workspaceID == "" {
		http.Error(w, "Workspace ID is required", http.StatusBadRequest)
		return
	}
	workspace, err := h.workspaceRepo.GetWorkspace(r.Context(), workspaceID)
	if err != nil {
		fmt.Println("Error getting workspace:", err)
		http.Error(w, "Failed to get workspace", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspace)
}

func (h *WorkspaceHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	workspaceId := r.PathValue("workspaceId")
	if workspaceId == "" {
		http.Error(w, "Workspace ID is required", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body for channel data
	var channelData struct {
		Name  string `json:"name"`
		Emoji string `json:"emoji"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&channelData); err != nil {
		// If parsing fails, use default values
		channelData.Name = ""
		channelData.Emoji = ""
	}

	channel, err := h.workspaceRepo.CreateChannel(r.Context(), workspaceId, channelData.Name, channelData.Emoji)
	if err != nil {
		http.Error(w, "Failed to create channel", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}