package middleware

import (
	"backend/pkg/utilities"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	RefreshTokenCookieName = "refresh_token"
)

func TokenAuthMiddleware(next http.Handler, redisClient *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error": "Missing or invalid Authorization header"}`, http.StatusUnauthorized)
			return
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Extract refresh token from cookies
		cookie, err := r.Cookie(RefreshTokenCookieName)
		if err != nil {
			http.Error(w, `{"error": "Missing or invalid refresh token cookie"}`, http.StatusUnauthorized)
			return
		}
		refreshToken := cookie.Value

		// Validate access token
		accessUserID, validAccess, err := utilities.ValidateAccessToken(r.Context(), redisClient, accessToken)
		if err != nil || !validAccess {
			http.Error(w, `{"error": "Invalid or expired access token"}`, http.StatusUnauthorized)
			return
		}

		// Validate refresh token
		refreshUserID, validRefresh, err := utilities.ValidateRefreshToken(r.Context(), redisClient, refreshToken)
		if err != nil || !validRefresh {
			http.Error(w, `{"error": "Invalid or expired refresh token"}`, http.StatusUnauthorized)
			return
		}

		// Ensure the user IDs match
		if accessUserID != refreshUserID {
			http.Error(w, `{"error": "Token mismatch"}`, http.StatusUnauthorized)
			return
		}

		// Add userID to the request context
		ctx := r.Context()
		ctx = utilities.WithUserID(ctx, accessUserID) // Assuming you have a helper to add userID to context
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
