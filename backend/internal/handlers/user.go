package handlers

import (
	"backend/internal/services"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"net/http"
	"os"
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
	router.Handle("/api/auth/user-login", middleware.RateLimitMiddleware(
		http.HandlerFunc(h.Login),
		h.limiter,
		time.Minute,
	))
	router.Handle("/user", middleware.RateLimitMiddleware(
		middleware.TokenAuthMiddleware(
			http.HandlerFunc(h.GetUser),
			h.store,
		),
		h.limiter,
		time.Minute,
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := utilities.GenerateAccessToken(r.Context(), h.store, user.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utilities.GenerateRefreshToken(r.Context(), h.store, user.ID)
	if err != nil {
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
