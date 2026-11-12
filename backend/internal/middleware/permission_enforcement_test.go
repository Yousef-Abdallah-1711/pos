package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/jwt"
)

type fakeEvaluator struct {
	level      string
	err        error
	gotUser    int64
	gotTenant  int64
	gotModule  int64
	callsCount int
}

func (f *fakeEvaluator) EffectiveLevel(_ context.Context, claims *jwt.Claims, moduleID int64) (string, error) {
	f.callsCount++
	f.gotUser, f.gotTenant, f.gotModule = int64(claims.UserID), int64(claims.TenantID), moduleID
	return f.level, f.err
}

func runPerm(t *testing.T, mw func(http.Handler) http.Handler, claims *jwt.Claims) (code int, ran bool, ctxModule int64, ctxLevel string) {
	t.Helper()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ran = true
		ctxModule, _ = r.Context().Value(RequiredModuleContextKey).(int64)
		ctxLevel, _ = r.Context().Value(RequiredPermissionContextKey).(string)
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
	if claims != nil {
		req = req.WithContext(context.WithValue(req.Context(), UserContextKey, claims))
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, ran, ctxModule, ctxLevel
}

var testClaims = &jwt.Claims{UserID: 7, TenantID: 3, RestaurantID: 5, Role: "employee"}

// SEC-D: in enforce mode the declared level is actually checked.
func TestPermission_EnforceLevels(t *testing.T) {
	cases := []struct {
		have, need string
		want       int
	}{
		{domain.PermissionNone, domain.PermissionRead, http.StatusForbidden},
		{domain.PermissionRead, domain.PermissionRead, http.StatusOK},
		{domain.PermissionRead, domain.PermissionWrite, http.StatusForbidden},
		{domain.PermissionWrite, domain.PermissionWrite, http.StatusOK},
		{domain.PermissionWrite, domain.PermissionDelete, http.StatusForbidden},
		{domain.PermissionDelete, domain.PermissionDelete, http.StatusOK},
		{domain.PermissionDelete, domain.PermissionAdmin, http.StatusForbidden},
		{domain.PermissionAdmin, domain.PermissionAdmin, http.StatusOK},
		{domain.PermissionAdmin, domain.PermissionRead, http.StatusOK},
		{domain.PermissionExport, domain.PermissionExport, http.StatusOK},
		{domain.PermissionRead, domain.PermissionExport, http.StatusForbidden},
	}
	for _, c := range cases {
		ev := &fakeEvaluator{level: c.have}
		mw := NewPermissionMiddleware(ev, PermissionModeEnforce)(4, c.need)
		code, ran, _, _ := runPerm(t, mw, testClaims)
		if code != c.want || ran != (c.want == http.StatusOK) {
			t.Errorf("have=%s need=%s: want %d, got %d (handler ran=%v)", c.have, c.need, c.want, code, ran)
		}
		if ev.gotModule != 4 || ev.gotTenant != 3 || ev.gotUser != 7 {
			t.Errorf("evaluator called with wrong identity/module: %+v", ev)
		}
	}
}

func TestPermission_UnauthenticatedIsRejected(t *testing.T) {
	ev := &fakeEvaluator{level: domain.PermissionAdmin}
	mw := NewPermissionMiddleware(ev, PermissionModeEnforce)(4, domain.PermissionRead)
	code, ran, _, _ := runPerm(t, mw, nil)
	if code != http.StatusUnauthorized || ran || ev.callsCount != 0 {
		t.Fatalf("want 401 without evaluating, got %d ran=%v calls=%d", code, ran, ev.callsCount)
	}
}

func TestKDSWebSocketReadPermissionIsRequired(t *testing.T) {
	for _, tc := range []struct {
		level string
		want  int
	}{
		{domain.PermissionNone, http.StatusForbidden},
		{domain.PermissionRead, http.StatusOK},
		{domain.PermissionWrite, http.StatusOK},
	} {
		evaluator := &fakeEvaluator{level: tc.level}
		middleware := NewPermissionMiddleware(evaluator, PermissionModeEnforce)(13, domain.PermissionRead)
		code, ran, module, level := runPerm(t, middleware, testClaims)
		if code != tc.want || ran != (tc.want == http.StatusOK) {
			t.Errorf("KDS level=%s: want status %d, got %d (handler ran=%v)", tc.level, tc.want, code, ran)
		}
		if evaluator.gotModule != 13 {
			t.Errorf("KDS permission context mismatch: evaluated module=%d, passed module=%d level=%q", evaluator.gotModule, module, level)
		}
		if tc.want == http.StatusOK && (module != 13 || level != domain.PermissionRead) {
			t.Errorf("allowed KDS request had wrong permission context: passed module=%d level=%q", module, level)
		}
	}
}

func TestPermission_EvaluatorFailureFailsClosedWhenEnforcing(t *testing.T) {
	ev := &fakeEvaluator{err: errors.New("db down")}
	mw := NewPermissionMiddleware(ev, PermissionModeEnforce)(4, domain.PermissionRead)
	code, ran, _, _ := runPerm(t, mw, testClaims)
	if code != http.StatusServiceUnavailable || ran {
		t.Fatalf("want 503, got %d ran=%v", code, ran)
	}
	mw = NewPermissionMiddleware(nil, PermissionModeEnforce)(4, domain.PermissionRead)
	code, ran, _, _ = runPerm(t, mw, testClaims)
	if code != http.StatusServiceUnavailable || ran {
		t.Fatalf("nil evaluator must fail closed, got %d ran=%v", code, ran)
	}
}

func TestPermission_AuditModeLogsButDoesNotBlock(t *testing.T) {
	ev := &fakeEvaluator{level: domain.PermissionNone}
	mw := NewPermissionMiddleware(ev, PermissionModeAudit)(4, domain.PermissionWrite)
	code, ran, mod, lvl := runPerm(t, mw, testClaims)
	if code != http.StatusOK || !ran {
		t.Fatalf("audit mode must not block, got %d ran=%v", code, ran)
	}
	if mod != 4 || lvl != domain.PermissionWrite {
		t.Fatalf("required module/level context keys must stay populated, got %d/%q", mod, lvl)
	}
	ev.err = errors.New("db down")
	if code, ran, _, _ = runPerm(t, mw, testClaims); code != http.StatusOK || !ran {
		t.Fatalf("audit mode must tolerate evaluator errors, got %d", code)
	}
}

func TestPermission_UnauthenticatedStillRejectedInAuditMode(t *testing.T) {
	mw := NewPermissionMiddleware(&fakeEvaluator{}, PermissionModeAudit)(4, domain.PermissionRead)
	if code, ran, _, _ := runPerm(t, mw, nil); code != http.StatusUnauthorized || ran {
		t.Fatalf("want 401, got %d ran=%v", code, ran)
	}
}

func TestParsePermissionMode(t *testing.T) {
	cases := []struct {
		explicit   string
		production bool
		want       PermissionMode
		wantErr    bool
	}{
		{"", true, PermissionModeEnforce, false},
		{"", false, PermissionModeAudit, false},
		{"enforce", false, PermissionModeEnforce, false},
		{" AUDIT ", false, PermissionModeAudit, false},
		{"audit", true, "", true}, // audit-only is not allowed in production
		{"off", false, "", true},
		{"nonsense", false, "", true},
	}
	for _, c := range cases {
		got, err := ParsePermissionMode(c.explicit, c.production)
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("ParsePermissionMode(%q,%v) = %q,%v want %q err=%v", c.explicit, c.production, got, err, c.want, c.wantErr)
		}
	}
}
