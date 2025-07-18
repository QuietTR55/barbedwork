package handlers

import (
	"backend/internal/services"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)
type UserAuthHandler struct {
	userService *services.UserService
	store       utilities.SessionStore
	limiter     ratelimiter.RateLimiter
}

func NewUserAuthHandler(userService *services.UserService, store utilities.SessionStore, limiter ratelimiter.RateLimiter) *UserAuthHandler {
	return &UserAuthHandler{
		userService: userService,
		store:       store,
		limiter:     limiter,
	}
}

func (h *UserAuthHandler) RegisterRoutes(router *http.ServeMux) {
	userAuthStack := []middleware.Middleware{
		middleware.RateLimitMiddleware(h.limiter, time.Minute, "user_auth"),
	}

	router.Handle("/api/auth/user-login", middleware.Chain(
		http.HandlerFunc(h.Login),
		userAuthStack...,
	))
}

func (h *UserAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		fmt.Println("Login error:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := utilities.GenerateAccessToken(r.Context(), h.store, user.Id.String())
	if err != nil {
		fmt.Println("Access token generation error:", err)
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utilities.GenerateRefreshToken(r.Context(), h.store, user.Id.String())
	if err != nil {
		fmt.Println("Refresh token generation error:", err)
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	maxAge := 3600 * 24 * 30
	sameSite := http.SameSiteLaxMode
	if os.Getenv("DEV") != "" {
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, &http.Cookie{
		Name:     utilities.RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   true,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken})
}