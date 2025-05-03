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
func GenerateAccessToken(ctx context.Context, store SessionStore, userId string) (string, error) {
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

	_, err = store.StoreSession(ctx, userId, int(time.Until(expiresAt).Seconds()))
	if err != nil {
		return "", errors.Join(ErrSessionStorage, err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new long-lived refresh token and stores it in Redis
func GenerateRefreshToken(ctx context.Context, store SessionStore, userId string) (string, error) {
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

	err = store.StoreRefreshToken(ctx, jti, userId, int(time.Until(expiresAt).Seconds()))
	if err != nil {
		return "", errors.Join(ErrSessionStorage, err)
	}

	return tokenString, nil
}

// ValidateAccessToken checks the access token and ensures the session is valid
// Returns the userId and an error if invalid.
func ValidateAccessToken(ctx context.Context, store SessionStore, tokenString string) (string, error) {
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

	val, err := store.GetSessionStatus(ctx, claims.Subject)
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
func ValidateRefreshToken(ctx context.Context, store SessionStore, tokenString string) (string, string, error) {
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

	userIdFromRedis, err := store.GetRefreshTokenUser(ctx, claims.Jti)
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
func RevokeUserSession(ctx context.Context, store SessionStore, userId string, jti string) error {
	err := store.RevokeSession(ctx, userId, jti)

	if err != nil && err != redis.Nil {
		return errors.Join(ErrSessionStorage, err)
	}

	return nil
}
