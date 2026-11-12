package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func withFixedKDSContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.TenantContextKey, int64(1))
		ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(2))
		ctx = context.WithValue(ctx, middleware.UserContextKey, &jwt.Claims{
			UserID:       7,
			TenantID:     1,
			RestaurantID: 2,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestKDSAnalyticsEndpointIntegrationSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}()

	repo := repository.NewKDSRepository(db)
	uc := usecase.NewKDSUseCase(repo)
	handler := NewKDSHandler(uc)

	var requiredModule interface{}
	var requiredLevel interface{}

	analyticsRoute := withFixedKDSContext(
		middleware.WithRequiredPermission(13, "ADMIN")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requiredModule = r.Context().Value(middleware.RequiredModuleContextKey)
				requiredLevel = r.Context().Value(middleware.RequiredPermissionContextKey)
				handler.GetAnalytics(w, r)
			}),
		),
	)

	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/kds/analytics", analyticsRoute)

	mock.ExpectQuery("average_prep_minutes").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"average_prep_minutes", "total_completed_orders"}).AddRow(9.3, int64(31)))
	mock.ExpectQuery("WITH days AS").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"date", "avg_prep_minutes", "completed_orders"}).AddRow("2026-02-20", 9.1, int64(5)))
	mock.ExpectQuery("WITH hours AS").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"hour", "order_count"}).AddRow(12, int64(11)))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/kds/analytics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	if requiredModule != int64(13) || requiredLevel != "ADMIN" {
		t.Fatalf("permission middleware context missing/incorrect: module=%v level=%v", requiredModule, requiredLevel)
	}

	var payload struct {
		Success bool `json:"success"`
		Data    struct {
			AveragePrepMinutes   float64 `json:"average_prep_minutes"`
			TotalCompletedOrders int64   `json:"total_completed_orders"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !payload.Success {
		t.Fatalf("expected success payload, got %+v", payload)
	}
	if payload.Data.AveragePrepMinutes != 9.3 || payload.Data.TotalCompletedOrders != 31 {
		t.Fatalf("unexpected analytics response: %+v", payload.Data)
	}
}

func TestKDSAnalyticsEndpointIntegrationRepoFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}()

	repo := repository.NewKDSRepository(db)
	uc := usecase.NewKDSUseCase(repo)
	handler := NewKDSHandler(uc)

	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/kds/analytics", withFixedKDSContext(http.HandlerFunc(handler.GetAnalytics)))

	mock.ExpectQuery("average_prep_minutes").
		WithArgs(int64(1), int64(2)).
		WillReturnError(context.DeadlineExceeded)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/kds/analytics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected 504, got %d body=%s", rec.Code, rec.Body.String())
	}

	var payload map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload["error"] != "KDS analytics request timed out" {
		t.Fatalf("unexpected error payload: %+v", payload)
	}
}
