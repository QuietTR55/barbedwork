package middleware

import (
	"backend/pkg/utilities"
	"net/http"
)

func PermissionMiddleware(requiredPermission string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, ok := r.Context().Value(utilities.UserIDKey).(string)
			if !ok || userId == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the user has the required permission
			hasPermission, err := utilities.CheckUserPermission(r.Context(), userId, requiredPermission)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}