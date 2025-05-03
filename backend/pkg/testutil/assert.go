package testutil

import (
	"strings"
	"testing"
)

func AssertEqual(t *testing.T, field string, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func AssertError(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("expected error, got nil")
	}
}

func AssertBody(t *testing.T, gotBody, wantBody string) {
	t.Helper()
	if strings.TrimSpace(gotBody) != strings.TrimSpace(wantBody) {
		t.Errorf("expected body %v, got %v", wantBody, gotBody)
	}
}

func AssertHeader(t *testing.T, headers map[string][]string, key, expectedValue string) {
	t.Helper()
	if val, ok := headers[key]; ok {
		if len(val) != 1 || val[0] != expectedValue {
			t.Errorf("unexpected header %s: got %q, want %q", key, val, expectedValue)
		}
	} else {
		t.Errorf("expected header %s, got none", key)
	}
}

func AssertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("expected status code %d, got %d", want, got)
	}
}

func AssertEqualBool(t *testing.T, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}
