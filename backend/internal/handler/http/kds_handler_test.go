package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func newKDSHandlerWithMock(t *testing.T) (*KDSHandler, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}

	repo := repository.NewKDSRepository(db)
	uc := usecase.NewKDSUseCase(repo)
	handler := NewKDSHandler(uc)

	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}

	return handler, mock, cleanup
}

func withKDSRequestContext(req *http.Request, tenantID, restaurantID, userID int64) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.TenantContextKey, tenantID)
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, restaurantID)
	if userID > 0 {
		ctx = context.WithValue(ctx, middleware.UserContextKey, &jwt.Claims{
			UserID:       int(userID),
			TenantID:     int(tenantID),
			RestaurantID: int(restaurantID),
		})
	}
	return req.WithContext(ctx)
}

func readErrorMessage(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()

	var payload map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if message, ok := payload["error"].(string); ok {
		return message
	}
	return ""
}

func TestKDSHandlerGetQueueMissingTenantContext(t *testing.T) {
	handler := NewKDSHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/kds/queue", nil)
	rec := httptest.NewRecorder()

	handler.GetQueue(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if !strings.Contains(readErrorMessage(t, rec), "Missing tenant context") {
		t.Fatalf("unexpected error body: %s", rec.Body.String())
	}
}

func TestKDSHandlerGetQueueInvalidStatusReturns400(t *testing.T) {
	handler, _, cleanup := newKDSHandlerWithMock(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/kds/queue?status=invalid_status", nil)
	req = withKDSRequestContext(req, 1, 2, 7)
	rec := httptest.NewRecorder()

	handler.GetQueue(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid status, got %d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(readErrorMessage(t, rec), "invalid kds status") {
		t.Fatalf("unexpected error body: %s", rec.Body.String())
	}
}

func TestKDSHandlerBumpOrderRejectsInvalidPathID(t *testing.T) {
	handler, _, cleanup := newKDSHandlerWithMock(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPut, "/api/v1/kds/orders/abc/bump", nil)
	req = withKDSRequestContext(req, 1, 2, 7)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	handler.BumpOrder(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id, got %d", rec.Code)
	}
	if !strings.Contains(readErrorMessage(t, rec), "Invalid KDS order ID") {
		t.Fatalf("unexpected error body: %s", rec.Body.String())
	}
}

func TestKDSHandlerBumpOrderSuccess(t *testing.T) {
	handler, mock, cleanup := newKDSHandlerWithMock(t)
	defer cleanup()

	receivedAt := time.Now().UTC().Add(-2 * time.Minute)
	startedAt := receivedAt.Add(30 * time.Second)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT status, order_id\\s+FROM kds_queue").
		WithArgs(int64(77), int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"status", "order_id"}).AddRow("new", int64(501)))
	mock.ExpectExec("UPDATE kds_queue\\s+SET").
		WithArgs("in_progress", int64(9), int64(77), int64(1), int64(2), "new").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE orders\\s+SET status = \\$1").
		WithArgs("preparing", int64(501), int64(1), int64(2), "confirmed").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO order_status_history").
		WithArgs(int64(501), "confirmed", "preparing", int64(9), "Synced from KDS bump").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("(?s)FROM kds_queue\\s+k.*WHERE k\\.id = \\$1 AND k\\.tenant_id = \\$2 AND k\\.restaurant_id = \\$3").
		WithArgs(int64(77), int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"tenant_id",
				"restaurant_id",
				"order_id",
				"order_number",
				"customer_name",
				"status",
				"priority",
				"station_id",
				"items",
				"notes",
				"received_at",
				"started_at",
				"completed_at",
				"bumped_by",
				"created_at",
				"updated_at",
				"table_id",
				"table_number",
				"table_name",
			}).AddRow(
				int64(77), int64(1), int64(2), int64(501), "ORD-501", "Alice", "in_progress", 0, "expo",
				[]byte("[]"), "", receivedAt, startedAt, nil, int64(9), receivedAt, startedAt, int64(31), int32(4), "Window",
			),
		)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/kds/orders/77/bump", nil)
	req = withKDSRequestContext(req, 1, 2, 9)
	req.SetPathValue("id", "77")
	rec := httptest.NewRecorder()

	handler.BumpOrder(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	var payload struct {
		Success bool `json:"success"`
		Data    struct {
			ID     int64  `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !payload.Success || payload.Data.ID != 77 || payload.Data.Status != "in_progress" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestKDSHandlerUpdateSettingsRejectsEmptyPayload(t *testing.T) {
	handler := NewKDSHandler(usecase.NewKDSUseCase(nil))

	req := httptest.NewRequest(http.MethodPut, "/api/v1/kds/settings", bytes.NewBufferString(`{}`))
	req = withKDSRequestContext(req, 1, 2, 0)
	rec := httptest.NewRecorder()

	handler.UpdateSettings(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(readErrorMessage(t, rec), "At least one settings field is required") {
		t.Fatalf("unexpected error body: %s", rec.Body.String())
	}
}

func TestKDSHandlerGetAnalyticsSuccess(t *testing.T) {
	handler, mock, cleanup := newKDSHandlerWithMock(t)
	defer cleanup()

	mock.ExpectQuery("average_prep_minutes").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"average_prep_minutes", "total_completed_orders"}).AddRow(7.4, int64(20)))
	mock.ExpectQuery("WITH days AS").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"date", "avg_prep_minutes", "completed_orders"}).AddRow("2026-02-20", 7.0, int64(4)))
	mock.ExpectQuery("WITH hours AS").
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"hour", "order_count"}).AddRow(11, int64(6)))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/kds/analytics", nil)
	req = withKDSRequestContext(req, 1, 2, 7)
	rec := httptest.NewRecorder()

	handler.GetAnalytics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
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
		t.Fatalf("expected success response: %+v", payload)
	}
	if payload.Data.AveragePrepMinutes != 7.4 || payload.Data.TotalCompletedOrders != 20 {
		t.Fatalf("unexpected analytics payload: %+v", payload.Data)
	}
}
