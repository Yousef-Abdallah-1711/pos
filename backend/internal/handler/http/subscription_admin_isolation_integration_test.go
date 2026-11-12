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

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestTenantTokenCannotExtendAnotherTenantsTrial(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	beforeA := readSubscriptionAdminSnapshot(t, db, fixture.tenantA)
	beforeB := readSubscriptionAdminSnapshot(t, db, fixture.tenantB)
	var auditBefore int
	if err := db.QueryRow(`SELECT COUNT(*) FROM super_admin_audit_logs WHERE target_id = $1`, fixture.tenantB).Scan(&auditBefore); err != nil {
		t.Fatalf("read target audit count before denied request: %v", err)
	}

	tenantTokens, err := jwt.NewTokenService("subscription-admin-tenant-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create tenant token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	tenantToken, err := tenantTokens.GenerateToken(1, fixture.tenantA, &restaurantID, "tenant-a@example.test", "owner")
	if err != nil {
		t.Fatalf("generate tenant token: %v", err)
	}
	superAdminTokens, err := jwt.NewSuperAdminTokenService("subscription-admin-platform-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create super-admin token service: %v", err)
	}
	handler := NewSubscriptionHandler(usecase.NewSubscriptionUseCase(
		repository.NewSubscriptionRepository(db), repository.NewSuperAdminAuditRepository(db),
	))
	secured := middleware.SuperAdminAuthMiddleware(superAdminTokens)(
		middleware.RequireSuperAdmin(repository.NewSuperAdminRepository(db))(
			http.HandlerFunc(handler.ExtendTenantTrial),
		),
	)

	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/super-admin/tenants/%d/extend-trial", fixture.tenantB), strings.NewReader(`{"days":30,"reason":"forged tenant request"}`))
	request.SetPathValue("id", fmt.Sprint(fixture.tenantB))
	request.Header.Set("Authorization", "Bearer "+tenantToken)
	request.Header.Set("X-Tenant-ID", fmt.Sprint(fixture.tenantB))
	request.Header.Set("X-Restaurant-ID", fmt.Sprint(fixture.restaurantB1))
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("tenant token reached platform subscription mutation: status=%d body=%s", response.Code, response.Body.String())
	}

	assertSubscriptionAdminSnapshot(t, db, fixture.tenantA, beforeA)
	assertSubscriptionAdminSnapshot(t, db, fixture.tenantB, beforeB)
	var auditAfter int
	if err := db.QueryRow(`SELECT COUNT(*) FROM super_admin_audit_logs WHERE target_id = $1`, fixture.tenantB).Scan(&auditAfter); err != nil {
		t.Fatalf("read target audit count after denied request: %v", err)
	}
	if auditAfter != auditBefore {
		t.Fatalf("denied tenant-token request changed target audit rows: before=%d after=%d", auditBefore, auditAfter)
	}
}

func TestSuperAdminTrialExtensionUpdatesOnlyTargetTenantAndAuditsActor(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	adminUserID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "subscription-super-admin")
	var superAdminID int64
	if err := db.QueryRow(`INSERT INTO super_admins (user_id, is_active) VALUES ($1, TRUE) RETURNING id`, adminUserID).Scan(&superAdminID); err != nil {
		t.Fatalf("insert super-admin fixture: %v", err)
	}
	baseEnd := time.Now().UTC().Truncate(time.Second)
	for _, tenantID := range []int{fixture.tenantA, fixture.tenantB} {
		if _, err := db.Exec(`UPDATE tenants SET subscription_status='expired', trial_end_date=$1, trial_days=4 WHERE id=$2`, baseEnd, tenantID); err != nil {
			t.Fatalf("seed tenant %d trial state: %v", tenantID, err)
		}
	}

	tokens, err := jwt.NewSuperAdminTokenService("subscription-admin-platform-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create super-admin token service: %v", err)
	}
	token, _, err := tokens.GenerateSuperAdminToken(adminUserID, "platform@example.test", superAdminID)
	if err != nil {
		t.Fatalf("generate super-admin token: %v", err)
	}
	handler := NewSubscriptionHandler(usecase.NewSubscriptionUseCase(
		repository.NewSubscriptionRepository(db), repository.NewSuperAdminAuditRepository(db),
	))
	secured := middleware.SuperAdminAuthMiddleware(tokens)(middleware.RequireSuperAdmin(repository.NewSuperAdminRepository(db))(
		http.HandlerFunc(handler.ExtendTenantTrial),
	))
	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/super-admin/tenants/%d/extend-trial", fixture.tenantA), strings.NewReader(`{"days":10,"reason":"support request"}`))
	request.SetPathValue("id", fmt.Sprint(fixture.tenantA))
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	secured.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("authorized trial extension returned %d: %s", response.Code, response.Body.String())
	}

	var targetDays, foreignDays int
	var targetEnd, foreignEnd time.Time
	var targetStatus, foreignStatus string
	if err := db.QueryRow(`SELECT trial_days, trial_end_date, subscription_status FROM tenants WHERE id=$1`, fixture.tenantA).Scan(&targetDays, &targetEnd, &targetStatus); err != nil {
		t.Fatalf("read target tenant after trial extension: %v", err)
	}
	if err := db.QueryRow(`SELECT trial_days, trial_end_date, subscription_status FROM tenants WHERE id=$1`, fixture.tenantB).Scan(&foreignDays, &foreignEnd, &foreignStatus); err != nil {
		t.Fatalf("read foreign tenant after trial extension: %v", err)
	}
	if targetDays != 14 || !targetEnd.Equal(baseEnd.AddDate(0, 0, 10)) || targetStatus != "trial" {
		t.Fatalf("target tenant trial state mismatch: days=%d end=%s status=%s", targetDays, targetEnd, targetStatus)
	}
	if foreignDays != 4 || !foreignEnd.Equal(baseEnd) || foreignStatus != "expired" {
		t.Fatalf("platform mutation changed foreign tenant: days=%d end=%s status=%s", foreignDays, foreignEnd, foreignStatus)
	}
	var historyCount int
	var action, reason string
	var actor int64
	if err := db.QueryRow(`SELECT COUNT(*), MAX(action), MAX(reason), MAX(performed_by) FROM tenant_subscription_history WHERE tenant_id=$1`, fixture.tenantA).Scan(&historyCount, &action, &reason, &actor); err != nil {
		t.Fatalf("read target subscription history: %v", err)
	}
	if historyCount != 1 || action != domain.ActionTrialExtended || reason != "support request" || actor != superAdminID {
		t.Fatalf("trial history mismatch: count=%d action=%q reason=%q actor=%d", historyCount, action, reason, actor)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM tenant_subscription_history WHERE tenant_id=$1`, fixture.tenantB).Scan(&historyCount); err != nil {
		t.Fatalf("read foreign subscription history: %v", err)
	}
	if historyCount != 0 {
		t.Fatalf("authorized target update created %d history rows for foreign tenant", historyCount)
	}
	var auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM super_admin_audit_logs WHERE super_admin_id=$1 AND action='UPDATE_SETTINGS' AND target_id=$2 AND details->>'operation'='EXTEND_TENANT_TRIAL'`, superAdminID, fixture.tenantA).Scan(&auditCount); err != nil {
		t.Fatalf("read platform audit record: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("authorized platform mutation audit count=%d, want 1", auditCount)
	}
}

type subscriptionAdminSnapshot struct {
	status       string
	trialDays    int
	trialEndDate sql.NullTime
}

func readSubscriptionAdminSnapshot(t *testing.T, db *sql.DB, tenantID int) subscriptionAdminSnapshot {
	t.Helper()
	var snapshot subscriptionAdminSnapshot
	err := db.QueryRow(`SELECT subscription_status, trial_days, trial_end_date FROM tenants WHERE id = $1`, tenantID).Scan(
		&snapshot.status, &snapshot.trialDays, &snapshot.trialEndDate,
	)
	if err != nil {
		t.Fatalf("read subscription snapshot for tenant %d: %v", tenantID, err)
	}
	return snapshot
}

func assertSubscriptionAdminSnapshot(t *testing.T, db *sql.DB, tenantID int, want subscriptionAdminSnapshot) {
	t.Helper()
	got := readSubscriptionAdminSnapshot(t, db, tenantID)
	if got.status != want.status || got.trialDays != want.trialDays || got.trialEndDate != want.trialEndDate {
		t.Fatalf("tenant %d subscription changed after denied admin request: got=%+v want=%+v", tenantID, got, want)
	}
}
