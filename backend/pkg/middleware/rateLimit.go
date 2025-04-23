package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const luaScript = `
local current = redis.call("INCR",KEYS[1])
if tonumber(current) == 1 then
	redis.call("EXPIRE",KEYS[1],ARGV[1])
end
return current
`

func RateLimitMiddleware(next http.Handler, rdb *redis.Client, maxRequests int, window time.Duration) http.Handler {
	sha, err := rdb.ScriptLoad(ctx, luaScript).Result()
	if err != nil {
		panic("failed to load lua script: " + err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userKey := getUserKey(r) // can be IP or user ID
		redisKey := fmt.Sprintf("rate_limit:%s", userKey)

		var count int
		count, err = rdb.EvalSha(ctx, sha, []string{redisKey}, int(window.Seconds())).Int()
		if err != nil {
			// fallback to Eval if script isn't cached
			if strings.Contains(err.Error(), "NOSCRIPT") {
				count, err = rdb.Eval(ctx, luaScript, []string{redisKey}, int(window.Seconds())).Int()
				if err != nil {
					http.Error(w, `{"error": "Rate limit error"}`, http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, `{"error": "Rate limit error"}`, http.StatusInternalServerError)
				return
			}
		}

		if count > maxRequests {
			http.Error(w, `{"error": "Too many requests. Please try again later."}`, http.StatusTooManyRequests)
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
