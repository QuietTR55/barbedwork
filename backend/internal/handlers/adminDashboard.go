package handlers

import (
	"backend/internal/repos"
	"backend/internal/services"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const AdminUserID = "admin"

type AdminDashboardHandler struct {
	SessionStore           utilities.SessionStore
	Limiter                ratelimiter.RateLimiter
	adminPanelPasswordHash []byte
	userService            *services.UserService
	workspaceRepo          *repos.WorkspaceRepo
}

func NewAdminDashboardHandler(SessionStore utilities.SessionStore, Limiter ratelimiter.RateLimiter, adminPanelPasswordHash []byte, userService *services.UserService, workspaceRepo *repos.WorkspaceRepo) *AdminDashboardHandler {
	return &AdminDashboardHandler{SessionStore: SessionStore, Limiter: Limiter, adminPanelPasswordHash: adminPanelPasswordHash, userService: userService, workspaceRepo: workspaceRepo}
}

func (h *AdminDashboardHandler) RegisterRoutes(router *http.ServeMux) {
	adminStack := []middleware.Middleware{
		middleware.RateLimitMiddleware(h.Limiter, time.Minute, "admin_dashboard"),
		middleware.TokenAuthMiddleware(h.SessionStore),
	}
	router.Handle("/api/admin/dashboard", middleware.Chain(
		http.HandlerFunc(h.GetDashboardData),
		adminStack...,
	))
	router.Handle("/api/admin/create-user", middleware.Chain(
		http.HandlerFunc(h.CreateNewUser),
		adminStack...,
	))

	router.Handle("/api/admin/users", middleware.Chain(
		http.HandlerFunc(h.GetAllUsers),
		adminStack...,
	))

	router.Handle("/api/admin/create-workspace", middleware.Chain(
		http.HandlerFunc(h.CreateWorkspace),
		adminStack...,
	))

	router.Handle("/api/admin/workspaces", middleware.Chain(
		http.HandlerFunc(h.GetWorkspaces),
		adminStack...,
	))

	router.Handle("/api/admin/workspaces/{workspaceId}", middleware.Chain(
		http.HandlerFunc(h.GetWorkspace),
		adminStack...,
	))

	router.Handle("/api/admin/workspaces/{workspaceId}/users/{userId}", middleware.Chain(
		http.HandlerFunc(h.AddUserToWorkspace),
		adminStack...,
	))
}

func (h *AdminDashboardHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}
	userIDFromContext, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userIDFromContext != AdminUserID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
}

func (h *AdminDashboardHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	userIDFromContext, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userIDFromContext != AdminUserID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	user, err := h.userService.CreateUser(r.Context(), credentials.Username, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AdminDashboardHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userId, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userId != AdminUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *AdminDashboardHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userId, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userId != AdminUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	imagePath, err := utilities.SaveImage(w, r)
	if err != nil {
		http.Error(w, "Unable to save the image", http.StatusBadRequest)
		return
	}
	workspaceName := r.FormValue("WorkspaceName")
	createdWorkspace, err := h.workspaceRepo.CreateWorkspace(r.Context(), workspaceName, imagePath)
	if err != nil {
		http.Error(w, "Unable to create workspace", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Workspace created successfully", "workspace": createdWorkspace})
}

func (h *AdminDashboardHandler) GetWorkspaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	userId, _ := utilities.GetUserID(ctx)
	if userId != AdminUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	workspaces, err := h.workspaceRepo.GetAllWorkspaces(ctx)
	if err != nil {
		http.Error(w, "Unable to get workspaces", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaces)
}

func (h *AdminDashboardHandler) GetWorkspace(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetWorkspace called")
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userId, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userId != AdminUserID {
		fmt.Println("Forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	workspaceID := r.PathValue("workspaceId")
	if workspaceID == "" {
		fmt.Println("Workspace ID is required")
		http.Error(w, "Workspace ID is required", http.StatusBadRequest)
		return
	}

	workspace, err := h.workspaceRepo.GetWorkspace(r.Context(), workspaceID)
	fmt.Println("Workspace ID:", workspaceID)
	fmt.Println("users", workspace.Users)
	if err != nil {
		fmt.Println("Unable to get workspace:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to get workspace","workspace": workspaceID})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspace)
}

func (h *AdminDashboardHandler) AddUserToWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userId, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userId != AdminUserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	workspaceID := r.PathValue("workspaceId")
	fmt.Println("Workspace ID:", workspaceID)
	if workspaceID == "" {
		http.Error(w, "Workspace ID is required", http.StatusBadRequest)
		return
	}

	userId = r.PathValue("userId")
	fmt.Println("User ID:", userId)
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	err := h.workspaceRepo.AddUserToWorkspace(r.Context(), userId, workspaceID)
	if err != nil {
		http.Error(w, "Unable to add user to workspace", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}