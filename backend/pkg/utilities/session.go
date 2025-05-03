package utilities

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


type SessionStore interface {
	GetSessionStatus(ctx context.Context, userId string) (string, error)
	GetRefreshTokenUser(ctx context.Context, jti string) (string, error)
	StoreSession(ctx context.Context, userId string, ttlSeconds int) (string, error)
	StoreRefreshToken(ctx context.Context, jti string, userId string, ttlSeconds int) error
	RevokeSession(ctx context.Context, userId string, jti string) error
}

type RedisSessionStore struct {
	Redis *redis.Client
}

func NewRedisSessionStore(redis *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{Redis: redis}
}

func (s *RedisSessionStore) GetSessionStatus(ctx context.Context, userId string) (string, error) {
	return s.Redis.Get(ctx, "session:"+userId).Result()
}

func (s *RedisSessionStore) GetRefreshTokenUser(ctx context.Context, jti string) (string, error) {
	return s.Redis.Get(ctx, "refresh:"+jti).Result()
}

func (s *RedisSessionStore) StoreSession(ctx context.Context, userId string, ttlSeconds int) (string, error) {
	return s.Redis.Set(ctx, "session:"+userId, "active", time.Duration(ttlSeconds)*time.Second).Result()
}

func (s *RedisSessionStore) StoreRefreshToken(ctx context.Context, jti string, userId string, ttlSeconds int) error {
	return s.Redis.Set(ctx, "refresh:"+jti, userId, time.Duration(ttlSeconds)*time.Second).Err()
}

func (s *RedisSessionStore) RevokeSession(ctx context.Context, userId string, jti string) error {
	pipe := s.Redis.Pipeline()
	pipe.Del(ctx, "session:"+userId)
	pipe.Del(ctx, "refresh:"+jti)
	_, err := pipe.Exec(ctx)
	return err
}
