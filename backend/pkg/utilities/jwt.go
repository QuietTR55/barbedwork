package utilities

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Custom errors for JWT validation
var (
	ErrTokenInvalid        = errors.New("token is invalid or expired")
	ErrSessionNotFound     = errors.New("session not found or inactive")
	ErrInvalidRefreshToken = errors.New("refresh token not found or mismatched")
	ErrRedisLookupFailed   = errors.New("failed to query session store")
	ErrTokenGeneration     = errors.New("failed to generate token")
	ErrSessionStorage      = errors.New("failed to store session data")
)

// AccessTokenClaims defines the structure of an access token
type AccessTokenClaims struct {
	jwt.RegisteredClaims
}

// RefreshTokenClaims defines the structure of a refresh token
type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	Jti string `json:"jti"`
}

// GenerateAccessToken creates a new short-lived access token
func GenerateAccessToken(ctx context.Context, redisClient *redis.Client, userId string) (string, error) {
	expiresAt := time.Now().Add(15 * time.Minute)

	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.Join(ErrTokenGeneration, err)
	}

	err = redisClient.Set(ctx, "session:"+userId, "active", time.Until(expiresAt)).Err()
	if err != nil {
		return "", errors.Join(ErrSessionStorage, err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new long-lived refresh token and stores it in Redis
func GenerateRefreshToken(ctx context.Context, redisClient *redis.Client, userId string) (string, error) {
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days
	jti := uuid.NewString()

	claims := RefreshTokenClaims{
		Jti: jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.Join(ErrTokenGeneration, err)
	}

	err = redisClient.Set(ctx, "refresh:"+jti, userId, time.Until(expiresAt)).Err()
	if err != nil {
		return "", errors.Join(ErrSessionStorage, err)
	}

	return tokenString, nil
}

// ValidateAccessToken checks the access token and ensures the session is valid
// Returns the userId and an error if invalid.
func ValidateAccessToken(ctx context.Context, redisClient *redis.Client, tokenString string) (string, error) {
	var claims AccessTokenClaims
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenMalformed) {
			return "", ErrTokenInvalid
		}
		return "", errors.Join(ErrTokenInvalid, err)
	}

	if !token.Valid {
		return "", ErrTokenInvalid
	}

	val, err := redisClient.Get(ctx, "session:"+claims.Subject).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrSessionNotFound
		}
		return "", errors.Join(ErrRedisLookupFailed, err)
	}

	if val != "active" {
		return "", ErrSessionNotFound
	}

	return claims.Subject, nil
}

// ValidateRefreshToken checks the refresh token and matches it against Redis
// Returns the userId and an error if invalid.
func ValidateRefreshToken(ctx context.Context, redisClient *redis.Client, tokenString string) (string, string, error) {
	var claims RefreshTokenClaims
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenMalformed) {
			return "", "", ErrTokenInvalid
		}
		return "", "", errors.Join(ErrTokenInvalid, err)
	}

	if !token.Valid || claims.Jti == "" {
		return "", "", ErrTokenInvalid
	}

	userIdFromRedis, err := redisClient.Get(ctx, "refresh:"+claims.Jti).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", ErrInvalidRefreshToken
		}
		return "", "", errors.Join(ErrRedisLookupFailed, err)
	}

	if userIdFromRedis != claims.Subject {
		return "", "", ErrInvalidRefreshToken
	}

	return claims.Subject, claims.Jti, nil
}

// RevokeUserSession deletes the access token session marker and the specific refresh token entry in Redis
// Requires the userId (for session marker) and the specific JTI of the refresh token being revoked.
func RevokeUserSession(ctx context.Context, redisClient *redis.Client, userId string, jti string) error {
	pipe := redisClient.Pipeline()
	pipe.Del(ctx, "session:"+userId)
	pipe.Del(ctx, "refresh:"+jti)
	_, err := pipe.Exec(ctx)

	if err != nil && err != redis.Nil {
		return errors.Join(ErrSessionStorage, err)
	}

	return nil
}

// RevokeAllUserRefreshTokens deletes all refresh tokens associated with a user ID.
// This requires iterating through keys, which can be inefficient on large datasets.
// Consider alternative structures if frequent mass revocation is needed.
// func RevokeAllUserRefreshTokens(ctx context.Context, redisClient *redis.Client, userId string) error {
//  	// Implementation would involve scanning keys matching "refresh:*" and checking the stored userId.
//  	// This is generally discouraged in production due to performance implications (SCAN).
//  	// A better approach might involve maintaining a set per user: e.g., user:<userId>:refresh_tokens = {jti1, jti2, ...}
//  	// Then revocation involves getting the set, deleting each "refresh:<jti>", and deleting the set.
//  	return errors.New("mass refresh token revocation strategy not implemented")
// }
