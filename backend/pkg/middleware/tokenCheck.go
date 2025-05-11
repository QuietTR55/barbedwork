package middleware

import (
	"backend/pkg/utilities"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const UserIDKey utilities.ContextKey = "userID"

func TokenAuthMiddleware(store utilities.SessionStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				fmt.Println("No Bearer token found")
				http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
				return
			}
			accessToken := strings.TrimPrefix(authHeader, "Bearer ")

			// Extract refresh token from cookies
			cookie, err := r.Cookie(utilities.RefreshTokenCookieName)
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

			// First try to validate the access token
			accessUserID, err := utilities.ValidateAccessToken(r.Context(), store, accessToken)
			if err == nil {
				// Access token is valid, proceed with the request
				ctx := r.Context()
				ctx = utilities.WithUserID(ctx, accessUserID)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			// Access token is invalid, try to validate refresh token
			refreshUserID, _, err := utilities.ValidateRefreshToken(r.Context(), store, refreshToken)
			if err != nil {
				fmt.Println("Refresh token invalid:", err)
				errorMsg := `{"error": "Session expired, please log in again"}`
				statusCode := http.StatusUnauthorized

				if errors.Is(err, utilities.ErrRedisLookupFailed) {
					errorMsg = `{"error": "Internal server error"}`
					statusCode = http.StatusInternalServerError
				}
				http.Error(w, errorMsg, statusCode)
				return
			}
			fmt.Println("Resfresh token: ", refreshToken)
			// Refresh token is valid, generate new access token
			fmt.Println("Generating new access token for user:", refreshUserID)
			newAccessToken, err := utilities.GenerateAccessToken(r.Context(), store, refreshUserID)
			if err != nil {
				fmt.Println("Failed to generate new access token:", err)
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}

			// Set the new access token in the response header
			w.Header().Set("X-New-Access-Token", newAccessToken)

			// Update the user ID for the rest of the middleware
			ctx := r.Context()
			ctx = utilities.WithUserID(ctx, refreshUserID)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
