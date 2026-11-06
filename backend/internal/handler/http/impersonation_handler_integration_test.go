package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

type impersonationHTTPFlow struct {
	db                *sql.DB
	admin             http.Handler
	handler           *ImpersonationHandler
	token             string
	targetUserID      int64
	superAdminID      int64
	store             *service.AuthSessionStore
	impersonationRepo *repository.ImpersonationRepository
	adminRepo         *repository.SuperAdminRepository
}

func TestImpersonationHTTPAuditSessionLifecycle(t *testing.T) {
	t.Setenv("ENV", "development")
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	adminUserID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "impersonation-admin")
	targetUserID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "impersonation-target")
	var superAdminID int64
	if err := db.QueryRow(`INSERT INTO super_admins (user_id, is_active) VALUES ($1, TRUE) RETURNING id`, adminUserID).Scan(&superAdminID); err != nil {
		t.Fatalf("insert super-admin fixture: %v", err)
	}
	tokens, err := jwt.NewSuperAdminTokenService("impersonation-http-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create super-admin token service: %v", err)
	}
	token, _, err := tokens.GenerateSuperAdminToken(adminUserID, "impersonation-admin@example.test", superAdminID)
	if err != nil {
		t.Fatalf("generate super-admin token: %v", err)
	}
	store := service.NewAuthSessionStore(service.NewInMemoryCache(), time.Hour)
	impersonationRepo := repository.NewImpersonationRepository(db)
	adminRepo := repository.NewSuperAdminRepository(db)
	uc := usecase.NewImpersonationUseCase(impersonationRepo, adminRepo, repository.NewSuperAdminAuditRepository(db))
	handler := NewImpersonationHandler(uc, store)
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("POST /api/v1/super-admin/impersonate/{userID}", handler.StartImpersonation)
	adminMux.HandleFunc("POST /api/v1/super-admin/impersonate/end", handler.EndImpersonation)
	adminMux.HandleFunc("GET /api/v1/super-admin/impersonate/status", handler.GetImpersonationStatus)
	securedAdminMux := middleware.SuperAdminAuthMiddleware(tokens)(middleware.RequireSuperAdmin(adminRepo)(adminMux))
	assertImpersonationSessionFlow(t, impersonationHTTPFlow{
		db: db, admin: securedAdminMux, handler: handler, token: token,
		targetUserID: targetUserID, superAdminID: superAdminID, store: store,
		impersonationRepo: impersonationRepo, adminRepo: adminRepo,
	})
	assertUnauditedImpersonationIsRevoked(t, db, securedAdminMux, targetUserID, adminUserID, superAdminID)
}

func assertImpersonationSessionFlow(t *testing.T, flow impersonationHTTPFlow) {
	t.Helper()
	start := authenticatedImpersonationRequest(http.MethodPost, fmt.Sprintf("/api/v1/super-admin/impersonate/%d", flow.targetUserID), flow.token, nil)
	startRecorder := httptest.NewRecorder()
	flow.admin.ServeHTTP(startRecorder, start)
	if startRecorder.Code != http.StatusOK {
		t.Fatalf("start impersonation returned %d: %s", startRecorder.Code, startRecorder.Body.String())
	}
	cookie := findImpersonationCookie(t, startRecorder.Result().Cookies())
	if !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Value == "" || strings.Contains(startRecorder.Body.String(), cookie.Value) {
		t.Fatalf("start response did not isolate the impersonation credential in an HttpOnly cookie")
	}
	assertImpersonatedUserCanReadOwnProfile(t, flow, cookie)
	assertImpersonationStatus(t, flow.admin, flow.token, flow.targetUserID)
	assertImpersonationValidation(t, flow.handler, cookie)
	assertImpersonationEnd(t, flow.db, flow.handler, cookie, flow.superAdminID)
}

func assertImpersonatedUserCanReadOwnProfile(t *testing.T, flow impersonationHTTPFlow, cookie *http.Cookie) {
	t.Helper()
	tokens, err := jwt.NewTokenService("impersonation-tenant-action-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create tenant token service: %v", err)
	}
	profileHandler := NewUserSettingsHandler(repository.NewUserSettingsRepository(flow.db), repository.NewUserRepository(flow.db))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/users/{userId}/profile", profileHandler.GetUserProfile)
	impersonationAuth := middleware.NewImpersonationAuthenticator(flow.impersonationRepo, flow.adminRepo)
	secured := middleware.SessionAuthMiddleware(tokens, flow.store, impersonationAuth)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: flow.db},
	)(mux))
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d/profile", flow.targetUserID), nil)
	request.AddCookie(cookie)
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), "impersonation-target@example.test") {
		t.Fatalf("impersonated user profile request returned %d: %s", recorder.Code, recorder.Body.String())
	}
}

func assertImpersonationStatus(t *testing.T, admin http.Handler, token string, targetUserID int64) {
	t.Helper()
	request := authenticatedImpersonationRequest(http.MethodGet, "/api/v1/super-admin/impersonate/status", token, nil)
	recorder := httptest.NewRecorder()
	admin.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), fmt.Sprintf(`"target_user_id":%d`, targetUserID)) {
		t.Fatalf("impersonation status mismatch: %d %s", recorder.Code, recorder.Body.String())
	}
}

func assertImpersonationValidation(t *testing.T, handler *ImpersonationHandler, cookie *http.Cookie) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/impersonation/validate", nil)
	request.AddCookie(cookie)
	recorder := httptest.NewRecorder()
	handler.ValidateImpersonationToken(recorder, request)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"valid":true`) || strings.Contains(recorder.Body.String(), "not-used-by-test") {
		t.Fatalf("impersonation validation mismatch or sensitive data leak: %d %s", recorder.Code, recorder.Body.String())
	}
}

func assertImpersonationEnd(t *testing.T, db *sql.DB, handler *ImpersonationHandler, cookie *http.Cookie, superAdminID int64) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/impersonation/end", nil)
	request.AddCookie(cookie)
	request.Header.Set("Origin", "http://localhost:3002")
	recorder := httptest.NewRecorder()
	middleware.RequireTrustedOrigin(http.HandlerFunc(handler.EndImpersonationByToken)).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("end impersonation returned %d: %s", recorder.Code, recorder.Body.String())
	}
	assertClearedImpersonationCookie(t, recorder.Result().Cookies())
	var active, auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM impersonation_sessions WHERE super_admin_id=$1 AND is_active`, superAdminID).Scan(&active); err != nil {
		t.Fatalf("count active sessions after end: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM super_admin_audit_logs WHERE super_admin_id=$1 AND action IN ('IMPERSONATE_START','IMPERSONATE_END')`, superAdminID).Scan(&auditCount); err != nil {
		t.Fatalf("count impersonation audit records: %v", err)
	}
	if active != 0 || auditCount != 2 {
		t.Fatalf("end did not revoke and audit session: active=%d audit_count=%d", active, auditCount)
	}
}

func assertUnauditedImpersonationIsRevoked(t *testing.T, db *sql.DB, admin http.Handler, targetUserID, adminUserID, superAdminID int64) {
	t.Helper()
	if _, err := db.Exec(`CREATE OR REPLACE FUNCTION reject_impersonation_start_audit() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.action = 'IMPERSONATE_START' THEN RAISE EXCEPTION 'audit unavailable'; END IF; RETURN NEW; END $$`); err != nil {
		t.Fatalf("create audit failure fixture function: %v", err)
	}
	if _, err := db.Exec(`CREATE TRIGGER reject_impersonation_start_audit BEFORE INSERT ON super_admin_audit_logs FOR EACH ROW EXECUTE FUNCTION reject_impersonation_start_audit()`); err != nil {
		t.Fatalf("create audit failure fixture trigger: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DROP TRIGGER IF EXISTS reject_impersonation_start_audit ON super_admin_audit_logs`)
		_, _ = db.Exec(`DROP FUNCTION IF EXISTS reject_impersonation_start_audit()`)
	})
	request := authenticatedImpersonationRequest(http.MethodPost, fmt.Sprintf("/api/v1/super-admin/impersonate/%d", targetUserID), testSuperAdminToken(t, adminUserID, superAdminID), nil)
	recorder := httptest.NewRecorder()
	admin.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError || len(recorder.Result().Cookies()) != 0 {
		t.Fatalf("unaudited start returned %d or issued a cookie: %s", recorder.Code, recorder.Body.String())
	}
	var active int
	if err := db.QueryRow(`SELECT COUNT(*) FROM impersonation_sessions WHERE super_admin_id=$1 AND is_active`, superAdminID).Scan(&active); err != nil {
		t.Fatalf("count sessions after failed audit: %v", err)
	}
	if active != 0 {
		t.Fatalf("failed audit left %d active impersonation sessions", active)
	}
}

func authenticatedImpersonationRequest(method, path, token string, cookie *http.Cookie) *http.Request {
	request := httptest.NewRequest(method, path, nil)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	if cookie != nil {
		request.AddCookie(cookie)
	}
	return request
}

func findImpersonationCookie(t *testing.T, cookies []*http.Cookie) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == service.ImpersonationSessionCookieName {
			return cookie
		}
	}
	t.Fatal("start response omitted impersonation session cookie")
	return nil
}

func assertClearedImpersonationCookie(t *testing.T, cookies []*http.Cookie) {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == service.ImpersonationSessionCookieName && cookie.MaxAge < 0 {
			return
		}
	}
	t.Fatal("end response did not clear impersonation cookie")
}

func testSuperAdminToken(t *testing.T, adminUserID, superAdminID int64) string {
	t.Helper()
	service, err := jwt.NewSuperAdminTokenService("impersonation-http-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create super-admin token service: %v", err)
	}
	token, _, err := service.GenerateSuperAdminToken(adminUserID, "admin@example.test", superAdminID)
	if err != nil {
		t.Fatalf("generate super-admin token: %v", err)
	}
	return token
}
