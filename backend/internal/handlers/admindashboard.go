package handlers

import (
	"encoding/json"
	"net/http" // Import net/http

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AdminDashboardHandler struct {
	redisClient            *redis.Client
	adminPanelPasswordHash []byte
}

func NewAdminDashboardHandler(redisClient *redis.Client, adminPanelPasswordHash []byte) *AdminDashboardHandler {
	return &AdminDashboardHandler{redisClient: redisClient, adminPanelPasswordHash: adminPanelPasswordHash}
}

const AdminUserID = "admin_user" // Define a constant for the admin user ID

// Example protected admin endpoint
func (h *AdminDashboardHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(gin.H{"error": "Method not allowed"})
		return
	}
	userIDFromContext, _ := r.Context().Value("userID").(string)
	if userIDFromContext != AdminUserID {
		// This shouldn't happen if middleware is correct, but as a safeguard
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(gin.H{"error": "Access denied"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gin.H{"message": "Welcome Admin! Dashboard data here."})
}
