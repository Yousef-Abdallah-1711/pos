package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func internalReq(t *testing.T, env string, production bool, headers map[string]string) int {
	t.Helper()
	h := RequireInternalToken(env, production)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health/metrics", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

// QA-001: operational endpoints are never open in production.
func TestInternalToken_ProductionWithoutSecretIsDisabled(t *testing.T) {
	t.Setenv("TEST_METRICS_TOKEN", "")
	if code := internalReq(t, "TEST_METRICS_TOKEN", true, nil); code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", code)
	}
}

func TestInternalToken_DevelopmentWithoutSecretIsOpen(t *testing.T) {
	t.Setenv("TEST_METRICS_TOKEN", "")
	if code := internalReq(t, "TEST_METRICS_TOKEN", false, nil); code != http.StatusOK {
		t.Fatalf("want 200, got %d", code)
	}
}

func TestInternalToken_WithSecretRequiresIt(t *testing.T) {
	t.Setenv("TEST_METRICS_TOKEN", "s3cret-value")
	cases := []struct {
		name    string
		headers map[string]string
		want    int
	}{
		{"none", nil, http.StatusUnauthorized},
		{"wrong header token", map[string]string{"X-Internal-Token": "nope"}, http.StatusUnauthorized},
		{"wrong bearer", map[string]string{"Authorization": "Bearer nope"}, http.StatusUnauthorized},
		{"empty bearer", map[string]string{"Authorization": "Bearer "}, http.StatusUnauthorized},
		{"correct header", map[string]string{"X-Internal-Token": "s3cret-value"}, http.StatusOK},
		{"correct bearer", map[string]string{"Authorization": "Bearer s3cret-value"}, http.StatusOK},
	}
	for _, production := range []bool{true, false} {
		for _, c := range cases {
			if got := internalReq(t, "TEST_METRICS_TOKEN", production, c.headers); got != c.want {
				t.Errorf("production=%v %s: want %d, got %d", production, c.name, c.want, got)
			}
		}
	}
}
