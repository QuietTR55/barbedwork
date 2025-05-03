package ratelimiter

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const luaScript = `
local current = redis.call("INCR",KEYS[1])
if tonumber(current) == 1 then
	redis.call("EXPIRE",KEYS[1],ARGV[1])
end
return current
`

type RateLimiter interface {
	Allow(ctx context.Context, key string, window time.Duration) (bool, error)
}

type RedisRateLimiter struct {
	rdb *redis.Client
	window time.Duration
	maxRequests int
	ScriptSha string
}

func NewRedisRateLimiter(rdb *redis.Client, window time.Duration, maxRequests int) (*RedisRateLimiter, error) {
	sha,err := rdb.ScriptLoad(context.Background(), luaScript).Result()
	if err != nil {
		return nil, err
	}

	return &RedisRateLimiter{
		rdb: rdb,
		window: window,
		maxRequests: maxRequests,
		ScriptSha: sha,
	}, nil
}

func (r *RedisRateLimiter) Allow(ctx context.Context, key string, window time.Duration) (bool, error) {
	count, err := r.rdb.EvalSha(ctx, r.ScriptSha, []string{key}, int(window.Seconds())).Int()
	if err != nil {
		if strings.Contains(err.Error(), "NOSCRIPT") {
			count, err = r.rdb.Eval(ctx, luaScript, []string{key}, int(window.Seconds())).Int()
			if err != nil {
				return false, err
			}
		} else {
				return false, err
		}
	}
	return count <= r.maxRequests, nil
}
