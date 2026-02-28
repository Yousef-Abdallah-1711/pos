package http

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	appjwt "pos-saas/internal/pkg/jwt"
)

func TestReportsRouteMiddlewareChain_RequiresAuthAndPropagatesRBACAndTenantContext(t *testing.T) {
	tokenService := newReportTestTokenService(t)

	wrapped := wrapReportsRouteForTest(
		tokenService,
		8,
		"READ",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotTenant := middleware.GetTenantID(r)
			gotRestaurant := middleware.GetRestaurantID(r)
			gotModule, _ := r.Context().Value(middleware.RequiredModuleContextKey).(int64)
			gotPermission, _ := r.Context().Value(middleware.RequiredPermissionContextKey).(string)

			if gotTenant != 8 || gotRestaurant != 8 {
				t.Fatalf("expected tenant/restaurant from JWT to be 8/8, got %d/%d", gotTenant, gotRestaurant)
			}
			if gotModule != 8 {
				t.Fatalf("expected required module 8, got %d", gotModule)
			}
			if gotPermission != "READ" {
				t.Fatalf("expected required permission READ, got %q", gotPermission)
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	t.Run("unauthorized without bearer token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/sales", nil)
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d body=%s", rr.Code, rr.Body.String())
		}
	})

	t.Run("jwt tenant context overrides conflicting headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/sales", nil)
		req.Header.Set("Authorization", "Bearer "+newReportTestJWT(t, tokenService, 7, 8, 8, "owner"))
		req.Header.Set("X-Tenant-ID", "999")
		req.Header.Set("X-Restaurant-ID", "999")
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestReportHandler_ExportRoute_UsesJWTTenantContextNotHeaders(t *testing.T) {
	tokenService := newReportTestTokenService(t)

	h := &ReportHandler{db: &sql.DB{}}
	h.salesReportDataBuilder = func(_, _ int64, _, _, _ time.Time, _ string, _ int64) (salesReportResponse, error) {
		return salesReportResponse{
			Filters:     map[string]interface{}{"date_from": "2026-02-01", "date_to": "2026-02-23"},
			Summary:     salesSummaryResponse{TotalRevenue: 10, TotalOrders: 1, AvgTicket: 10},
			GeneratedAt: "2026-02-23T00:00:00Z",
		}, nil
	}

	var auditCalls []reportExportAuditRecord
	h.reportExportAuditLogger = func(_ context.Context, rec reportExportAuditRecord) error {
		auditCalls = append(auditCalls, rec)
		return nil
	}

	wrapped := wrapReportsRouteForTest(
		tokenService,
		8,
		"EXPORT",
		http.HandlerFunc(h.ExportReport),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/reports/export",
		strings.NewReader(`{"report_type":"sales","format":"csv"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+newReportTestJWT(t, tokenService, 7, 8, 8, "owner"))
	req.Header.Set("X-Tenant-ID", "999")
	req.Header.Set("X-Restaurant-ID", "999")
	req.Header.Set("X-Forwarded-For", "198.51.100.20")
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}
	if len(auditCalls) != 1 {
		t.Fatalf("expected 1 audit call, got %d", len(auditCalls))
	}
	if auditCalls[0].TenantID != 8 || auditCalls[0].RestaurantID != 8 || auditCalls[0].UserID != 7 {
		t.Fatalf("expected audit IDs from JWT context (8/8/7), got %+v", auditCalls[0])
	}
}

func TestReportsRouteMiddlewareChain_RejectsMissingRestaurantContext(t *testing.T) {
	tokenService := newReportTestTokenService(t)

	// Use the real export handler to exercise tenant verification after middleware.
	h := &ReportHandler{db: &sql.DB{}}
	wrapped := wrapReportsRouteForTest(
		tokenService,
		8,
		"EXPORT",
		http.HandlerFunc(h.ExportReport),
	)

	// Generate a token with tenant but no restaurant to ensure TenantContextMiddleware rejects.
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/export", strings.NewReader(`{"report_type":"sales","format":"csv"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+newReportTestJWTNoRestaurant(t, tokenService, 7, 8, "owner"))
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing restaurant tenant context, got %d body=%s", rr.Code, rr.Body.String())
	}
}

func wrapReportsRouteForTest(tokenService *appjwt.TokenService, moduleID int64, permission string, next http.Handler) http.Handler {
	wrapped := middleware.WithRequiredPermission(moduleID, permission)(next)
	wrapped = middleware.TenantContextMiddleware(wrapped)
	wrapped = middleware.AuthMiddleware(tokenService)(wrapped)
	return wrapped
}

func newReportTestTokenService(t *testing.T) *appjwt.TokenService {
	t.Helper()
	svc, err := appjwt.NewTokenService("test-reports-jwt-secret", "24h")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	return svc
}

func newReportTestJWT(t *testing.T, tokenService *appjwt.TokenService, userID, tenantID, restaurantID int, role string) string {
	t.Helper()
	token, err := tokenService.GenerateToken(userID, tenantID, &restaurantID, "reports-test@example.com", role)
	if err != nil {
		t.Fatalf("generate jwt: %v", err)
	}
	return token
}

func newReportTestJWTNoRestaurant(t *testing.T, tokenService *appjwt.TokenService, userID, tenantID int, role string) string {
	t.Helper()
	token, err := tokenService.GenerateToken(userID, tenantID, nil, "reports-test@example.com", role)
	if err != nil {
		t.Fatalf("generate jwt: %v", err)
	}
	return token
}
