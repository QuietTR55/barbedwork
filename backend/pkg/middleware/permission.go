package middleware

import (
	"backend/pkg/utilities"
	"net/http"
)

func PermissionMiddleware(permissionChecker *utilities.PermissionChecker, requiredPermissions ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, ok := r.Context().Value(utilities.UserIDKey).(string)
			if !ok || userId == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Admin users bypass all permission checks
			if userId == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if the user has any of the required permissions
			hasPermission := false
			for _, permission := range requiredPermissions {
				userHasPermission, err := permissionChecker.CheckUserPermission(r.Context(), userId, permission)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				if userHasPermission {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}