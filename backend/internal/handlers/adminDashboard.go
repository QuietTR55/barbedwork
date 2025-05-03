package handlers

import (
	"backend/internal/services"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"fmt"
	"net/http" // Import net/http
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminDashboardHandler struct {
	SessionStore           utilities.SessionStore
	Limiter                ratelimiter.RateLimiter
	adminPanelPasswordHash []byte
	userService            *services.UserService
}

func NewAdminDashboardHandler(SessionStore utilities.SessionStore, Limiter ratelimiter.RateLimiter, adminPanelPasswordHash []byte, userService *services.UserService) *AdminDashboardHandler {
	return &AdminDashboardHandler{SessionStore: SessionStore, Limiter: Limiter, adminPanelPasswordHash: adminPanelPasswordHash, userService: userService}
}

func (h *AdminDashboardHandler) RegisterRoutes(router *http.ServeMux) {
	router.Handle("/api/admin/dashboard", middleware.RateLimitMiddleware(
		middleware.TokenAuthMiddleware(
			http.HandlerFunc(h.GetDashboardData),
			h.SessionStore,
		),
		h.Limiter,
		time.Minute,
	))
	router.Handle("/api/admin/create-user", middleware.RateLimitMiddleware(
		middleware.TokenAuthMiddleware(
			http.HandlerFunc(h.CreateNewUser),
			h.SessionStore,
		),
		h.Limiter,
		time.Minute,
	))
	router.Handle("/api/admin/users", middleware.RateLimitMiddleware(
		middleware.TokenAuthMiddleware(
			http.HandlerFunc(h.GetAllUsers),
			h.SessionStore,
		),
		h.Limiter,
		time.Minute,
	))
}

const AdminUserID = "admin" // Define a constant for the admin user ID

// Example protected admin endpoint
func (h *AdminDashboardHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(gin.H{"error": "Method not allowed"})
		return
	}
	userIDFromContext, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userIDFromContext != AdminUserID {
		// This shouldn't happen if middleware is correct, but as a safeguard
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(gin.H{"error": "Access denied"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gin.H{"message": "Welcome Admin! Dashboard data here."})
}

func (h *AdminDashboardHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(gin.H{"error": "Method not allowed"})
		return
	}

	userIDFromContext, _ := r.Context().Value(utilities.UserIDKey).(string)
	fmt.Println("User ID:", userIDFromContext)
	if userIDFromContext != AdminUserID {
		fmt.Println("User is not admin")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(gin.H{"error": "Access denied"})
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(gin.H{"error": "Invalid request body"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(gin.H{"error": "Failed to hash password"})
		return
	}

	//insert into users table
	user, err := h.userService.CreateUser(r.Context(), credentials.Username, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(gin.H{"error": "Failed to create user"})
		fmt.Println("Error creating user:", err)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gin.H{"message": "User created successfully", "user": user})
}

func (h *AdminDashboardHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(gin.H{"error": "Method not allowed"})
		return
	}

	userId, _ := r.Context().Value(utilities.UserIDKey).(string)
	if userId != AdminUserID {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(gin.H{"error": "Access denied"})
		return
	}
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(gin.H{"error": "Failed to get users"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gin.H{"users": users})
}