package middleware

import (
	customErrors "backend/errors"
	"backend/pkg/testutil"
	"backend/pkg/utilities"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeStore struct {
	failRedis          bool
	validAccessTokens  map[string]string
	validRefreshTokens map[string]string
	revokedSessions    map[string]bool
}

func (s *fakeStore) GetSessionStatus(ctx context.Context, userId string) (string, error) {
	if s.failRedis {
		return "", errors.New(customErrors.ErrRedisLookupFailed)
	}
	return "active", nil
}

func (s *fakeStore) GetRefreshTokenUser(ctx context.Context, jti string) (string, error) {
	return "user123", nil
}

func (s *fakeStore) StoreSession(ctx context.Context, userId string, ttlSeconds int) (string, error) {
	return "session123", nil
}

func (s *fakeStore) StoreRefreshToken(ctx context.Context, jti string, userId string, ttlSeconds int) error {
	return nil
}

func (s *fakeStore) RevokeSession(ctx context.Context, userId string, jti string) error {
	return nil
}

func TestTokenAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name            string
		expectedStatus  int
		expectNewAccess bool
		fakeStore       *fakeStore
		setupRequest    func(req *http.Request, store *fakeStore)
	}{
		{
			name:            "Valid access token",
			expectedStatus:  http.StatusOK,
			expectNewAccess: false,
			fakeStore: &fakeStore{
				validAccessTokens:  make(map[string]string),
				validRefreshTokens: make(map[string]string),
				revokedSessions:    make(map[string]bool),
			},
			setupRequest: func(req *http.Request, store *fakeStore) {
				accessToken, err := utilities.GenerateAccessToken(context.Background(), store, "user123")
				if err != nil {
					t.Fatalf("Failed to generate access token: %v", err)
				}
				refreshToken, err := utilities.GenerateRefreshToken(context.Background(), store, "user123")
				if err != nil {
					t.Fatalf("Failed to generate refresh token: %v", err)
				}
				req.Header.Set("Authorization", "Bearer "+accessToken)
				req.AddCookie(&http.Cookie{Name: utilities.RefreshTokenCookieName, Value: refreshToken})
			},
		},
		{
			name:            "invalid access token with valid refresh token",
			expectedStatus:  http.StatusOK,
			expectNewAccess: true,
			fakeStore: &fakeStore{
				validAccessTokens:  make(map[string]string),
				validRefreshTokens: make(map[string]string),
				revokedSessions:    make(map[string]bool),
			},
			setupRequest: func(req *http.Request, store *fakeStore) {
				refreshToken, err := utilities.GenerateRefreshToken(context.Background(), store, "user123")
				if err != nil {
					t.Fatalf("Failed to generate refresh token: %v", err)
				}
				req.Header.Set("Authorization", "Bearer "+"the not valid access token")
				req.AddCookie(&http.Cookie{Name: utilities.RefreshTokenCookieName, Value: refreshToken})
			},
		},
		{
			name:            "invalid access token with invalid refresh token",
			expectedStatus:  http.StatusUnauthorized,
			expectNewAccess: false,
			fakeStore: &fakeStore{
				validAccessTokens:  make(map[string]string),
				validRefreshTokens: make(map[string]string),
				revokedSessions:    make(map[string]bool),
			},
			setupRequest: func(req *http.Request, store *fakeStore) {
				req.Header.Set("Authorization", "Bearer "+"the not valid access token")
				req.AddCookie(&http.Cookie{Name: utilities.RefreshTokenCookieName, Value: "the not valid refresh token"})
			},
		},
		{
			name:            "Valid access token with invalid refresh token",
			expectedStatus:  http.StatusUnauthorized,
			expectNewAccess: false,
			fakeStore: &fakeStore{
				validAccessTokens:  make(map[string]string),
				validRefreshTokens: make(map[string]string),
				revokedSessions:    make(map[string]bool),
			},
			setupRequest: func(req *http.Request, store *fakeStore) {
				accessToken, err := utilities.GenerateAccessToken(context.Background(), store, "user123")
				if err != nil {
					t.Fatalf("Failed to generate access token: %v", err)
				}
				req.Header.Set("Authorization", "Bearer "+accessToken)
				req.AddCookie(&http.Cookie{Name: utilities.RefreshTokenCookieName, Value: "the not valid refresh token"})
			},
		},
		{
			name:            "invalid request",
			expectedStatus:  http.StatusBadRequest,
			expectNewAccess: false,
			fakeStore: &fakeStore{
				validAccessTokens:  make(map[string]string),
				validRefreshTokens: make(map[string]string),
				revokedSessions:    make(map[string]bool),
			},
			setupRequest: func(req *http.Request, store *fakeStore) {
				req.Header.Set("Authorization", "header with no Bearer")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			handler := TokenAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.expectedStatus)
			}), testCase.fakeStore)

			req := httptest.NewRequest("GET", "/test", nil)
			testCase.setupRequest(req, testCase.fakeStore)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			testutil.AssertStatusCode(t, recorder.Code, testCase.expectedStatus)
			returnedAccessToken := recorder.Header().Get("X-New-Access-Token")
			hasReturnedAccessToken := returnedAccessToken != ""
			testutil.AssertEqualBool(t, hasReturnedAccessToken, testCase.expectNewAccess)

		})
	}
}
