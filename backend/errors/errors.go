package errors

const (
	ErrRedisLookupFailed = "redis lookup failed"
	ErrRedisSetFailed = "redis set failed"
	ErrRedisDeleteFailed = "redis delete failed"
	ErrRedisPipelineFailed = "redis pipeline failed"
	ErrTokenGeneration = "token generation failed"
	ErrSessionStorage = "session storage failed"
	ErrTokenValidation = "token validation failed"
	ErrTokenRevocation = "token revocation failed"
	ErrTokenExpiration = "token expired"
	ErrTokenInvalid = "token invalid"
	ErrTokenMissing = "token missing"

)