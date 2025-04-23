package utilities

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func GenerateRefreshToken(ctx context.Context, redisClient *redis.Client, userId string) (string, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 30) // 30 days

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     expiresAt.Unix(), // Use standard 'exp' claim
		"type":    "refresh",
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	// Store refresh token reference/validity marker in Redis
	// Storing the full token might be redundant if it's in the cookie,
	// but useful for revocation checks. Adjust key/value as needed.
	err = redisClient.Set(ctx, "refresh:"+userId, tokenString, time.Until(expiresAt)).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateAccessToken(ctx context.Context, redisClient *redis.Client, userId string) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 15) // 15 minutes

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     expiresAt.Unix(), // Use standard 'exp' claim
		"type":    "access",
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	// Store an active session marker in Redis (optional but good for quick checks/revocation)
	err = redisClient.Set(ctx, "session:"+userId, "active", time.Until(expiresAt)).Err()
	if err != nil {
		// Log error, but maybe don't fail token generation just for this? Depends on requirements.
		// return "", err
	}

	return tokenString, nil
}

func ValidateAccessToken(ctx context.Context, redisClient *redis.Client, tokenString string) (string, bool, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Add algorithm validation if needed
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", false, err // Handles expired signature, invalid signature etc.
	}

	if !token.Valid {
		return "", false, nil // Token is invalid for other reasons
	}

	// Extract user ID
	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", false, nil // Invalid claim format
	}

	// Check token type (optional but good practice)
	tokenType, _ := claims["type"].(string)
	if tokenType != "access" {
		return "", false, nil // Wrong token type
	}

	// Check Redis for active session (optional, depends if you store session markers)
	val, err := redisClient.Get(ctx, "session:"+userId).Result()
	if err == redis.Nil {
		return "", false, nil // Session explicitly revoked or expired in Redis
	} else if err != nil {
		return "", false, err // Redis error
	}
	if val != "active" {
		return "", false, nil // Session marked inactive
	}

	return userId, true, nil
}

func ValidateRefreshToken(ctx context.Context, redisClient *redis.Client, tokenString string) (string, bool, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", false, err // Handles expired signature, invalid signature etc.
	}

	if !token.Valid {
		return "", false, nil
	}

	// Extract user ID
	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", false, nil
	}

	// Check token type
	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return "", false, nil // Wrong token type
	}

	// Check if stored token in Redis matches the provided one (for revocation)
	storedToken, err := redisClient.Get(ctx, "refresh:"+userId).Result()
	if err == redis.Nil {
		return "", false, nil // Token not found in Redis (revoked or never issued?)
	} else if err != nil {
		return "", false, err // Redis error
	}

	if storedToken != tokenString {
		return "", false, nil // Provided token doesn't match the one stored in Redis (likely revoked)
	}

	// Expiry is already checked by jwt.ParseWithClaims if 'exp' claim is present and valid

	return userId, true, nil
}

func RevokeUserSession(ctx context.Context, redisClient *redis.Client, userId string) error {
	// Remove both access and refresh tokens from Redis
	err := redisClient.Del(ctx, "session:"+userId).Err()
	if err != nil {
		return err
	}

	err = redisClient.Del(ctx, "refresh:"+userId).Err()
	if err != nil {
		return err
	}

	return nil
}
