package middleware

import (
	"backend/pkg/ratelimiter"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const luaScript = `
local current = redis.call("INCR",KEYS[1])
if tonumber(current) == 1 then
	redis.call("EXPIRE",KEYS[1],ARGV[1])
end
return current
`

func RateLimitMiddleware(next http.Handler, limiter ratelimiter.RateLimiter, window time.Duration) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userKey := getUserKey(r) // can be IP or user ID
		redisKey := fmt.Sprintf("rate_limit:%s", userKey)

		allowed, err := limiter.Allow(r.Context(), redisKey, window)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Rate limit error"}`))
			return
		}

		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Too many requests. Please try again later."}`))
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func getUserKey(r *http.Request) string {
	// Use the client's IP address as the key
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}
	return ip
}
