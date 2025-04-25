package middleware

import (
	"backend/pkg/utilities"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	RefreshTokenCookieName = "refreshToken"
)

const UserIDKey utilities.ContextKey = "userID"

func TokenAuthMiddleware(next http.Handler, redisClient *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("No Bearer token found")
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Extract refresh token from cookies
		cookie, err := r.Cookie(RefreshTokenCookieName)
		if err != nil {
			fmt.Println("No refresh token found")
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			} else {
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			}
			return
		}
		refreshToken := cookie.Value
		
		refreshUserID, _, err := utilities.ValidateRefreshToken(r.Context(), redisClient, refreshToken)
		var refreshTokenValid bool = err == nil

		if !refreshTokenValid {
			fmt.Println("Refresh token invalid")
			errorMsg := `{"error": "Session expired, please log in again"}`
			statusCode := http.StatusUnauthorized
			
			if errors.Is(err, utilities.ErrRedisLookupFailed) {
				errorMsg = `{"error": "Internal server error"}`
				statusCode = http.StatusInternalServerError
			}
			http.Error(w, errorMsg, statusCode)
			return
		}

		// First try to validate the access token
		accessUserID, err := utilities.ValidateAccessToken(r.Context(), redisClient, accessToken)
		var accessTokenValid bool = err == nil

		// If access token is invalid, check the refresh token
		if !accessTokenValid {
			fmt.Println("Access token invalid, checking refresh token")
			
			// At this point, refresh token is valid but access token is not
			fmt.Println("Generating new access token for user:", refreshUserID)
			
			// Generate a new access token
			newAccessToken, err := utilities.GenerateAccessToken(r.Context(), redisClient, refreshUserID)
			if err != nil {
				fmt.Println("Failed to generate new access token:", err)
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}
			
			// Set the new access token in the response header
			// The client-side code will need to extract this and update its stored token
			w.Header().Set("X-New-Access-Token", newAccessToken)
			
			// Update the user ID for the rest of the middleware
			accessUserID = refreshUserID
		}

		// At this point, we have a valid user ID (either from the original access token or the refresh token)
		// Add userID to the request context
		ctx := r.Context()
		ctx = utilities.WithUserID(ctx, accessUserID)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
