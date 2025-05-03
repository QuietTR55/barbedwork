package handlers

import (
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const RefreshTokenCookieName = "refreshToken"

type AdminAuthHandler struct {
	adminPanelPasswordHash []byte
	SessionStore           utilities.SessionStore
	Limiter                ratelimiter.RateLimiter
}

// NewAdminAuthHandler creates a new instance of AdminAuthHandler
func NewAdminAuthHandler(adminPanelPasswordHash []byte, sessionStore utilities.SessionStore, limiter ratelimiter.RateLimiter) *AdminAuthHandler {
	return &AdminAuthHandler{adminPanelPasswordHash: adminPanelPasswordHash, SessionStore: sessionStore, Limiter: limiter}
}

func (h *AdminAuthHandler) RegisterRoutes(router *http.ServeMux) {
	router.Handle("/api/auth/admin/login", middleware.RateLimitMiddleware(
		http.HandlerFunc(h.Login),
		h.Limiter,
		time.Minute,
	))
}

// Login handles admin login requests
func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method: " + r.Method)
		fmt.Println("Request URL: " + r.URL.String())
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SecretKey string `json:"secretKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Compare the hash from the container with the provided password
	err := bcrypt.CompareHashAndPassword(h.adminPanelPasswordHash, []byte(req.SecretKey))
	if err != nil {
		http.Error(w, `{"error": "Invalid secret password"}`, http.StatusUnauthorized)
		return
	}

	// --- Password is correct, generate tokens for the admin user ---
	userId := AdminUserID // Use the predefined admin user ID

	accessToken, err := utilities.GenerateAccessToken(r.Context(), h.SessionStore, userId)
	if err != nil {
		http.Error(w, `{"error": "Could not generate access token"}`, http.StatusInternalServerError)
		return
	}

	refreshToken, err := utilities.GenerateRefreshToken(r.Context(), h.SessionStore, userId)
	if err != nil {
		http.Error(w, `{"error": "Could not generate refresh token"}`, http.StatusInternalServerError)
		return
	}

	// Set the refresh token as a cookie
	maxAge := 3600 * 24 * 30
	sameSite := http.SameSiteLaxMode
	if os.Getenv("DEV") != "" {
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: sameSite,
	})

	// Return the access token in the response body
	json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken})
	w.WriteHeader(http.StatusOK)
}
