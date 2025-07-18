package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeLimiter struct {
	allowed bool
	err     error
}

func (l *fakeLimiter) Allow(ctx context.Context, key string, window time.Duration) (bool, error) {
	return l.allowed, l.err
}

func assertEqual(t *testing.T, field string, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("%s: expected %v, got %v", field, want, got)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name                     string
		allowed                  bool
		wantStatus               int
		forceError               bool
		expectFinalHandlerCalled bool
		wantBody                 string
		wantContentType          string
	}{
		{
			name:                     "allowed request",
			allowed:                  true,
			wantStatus:               http.StatusOK,
			forceError:               false,
			expectFinalHandlerCalled: true,
			wantBody:                 "",
			wantContentType:          "",
		},
		{
			name:                     "blocked request",
			allowed:                  false,
			wantStatus:               http.StatusTooManyRequests,
			forceError:               false,
			expectFinalHandlerCalled: false,
			wantBody:                 `{"error": "Too many requests. Please try again later."}`,
			wantContentType:          "application/json",
		},
		{
			name:                     "internal error in rate limiter",
			allowed:                  true, // doesn't matter because forceError=true
			wantStatus:               http.StatusInternalServerError,
			forceError:               true,
			expectFinalHandlerCalled: false,
			wantBody:                 `{"error": "Rate limit error"}`,
			wantContentType:          "application/json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var err error
			if test.forceError {
				err = fmt.Errorf("simulated error")
			}

			finalHandlerCalled := false
			finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				finalHandlerCalled = true
				w.WriteHeader(test.wantStatus)
			})

			middlewareStack := []Middleware{
				RateLimitMiddleware(&fakeLimiter{
					allowed: test.allowed,
					err:     err,
				}, time.Minute, "test_context"),
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			rec := httptest.NewRecorder()

			Chain(finalHandler, middlewareStack...).ServeHTTP(rec, req)

			assertEqual(t, "status code", rec.Code, test.wantStatus)

			if test.wantBody != "" {
				body := strings.TrimSpace(rec.Body.String())
				assertEqual(t, "body", body, test.wantBody)

				contentType := rec.Header().Get("Content-Type")
				assertEqual(t, "content type", contentType, test.wantContentType)
			}

			if finalHandlerCalled != test.expectFinalHandlerCalled {
				t.Errorf("final handler called = %v; want %v", finalHandlerCalled, test.expectFinalHandlerCalled)
			}
		})
	}
}
