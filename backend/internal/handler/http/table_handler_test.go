package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func newTableHandlerWithMock(t *testing.T) (*TableHandler, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	repo := repository.NewTableRepository(db)
	uc := usecase.NewTableUseCase(repo)
	handler := NewTableHandler(uc)

	return handler, mock, func() {
		_ = db.Close()
	}
}

func withTenantContext(r *http.Request, tenantID, restaurantID int64) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.TenantContextKey, tenantID)
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, restaurantID)
	return r.WithContext(ctx)
}

func TestTableHandler_CreateTableMissingTenantContext(t *testing.T) {
	handler := NewTableHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tables", bytes.NewBufferString(`{}`))
	res := httptest.NewRecorder()

	handler.CreateTable(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing tenant context, got %d", res.Code)
	}
}

func TestTableHandler_CreateTableSuccess(t *testing.T) {
	handler, mock, cleanup := newTableHandlerWithMock(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO restaurant_tables")).
		WithArgs(
			int64(1),
			int64(2),
			11,
			sqlmock.AnyArg(),
			4,
			"indoor",
			"available",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			0.0,
			0.0,
			"square",
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(int64(11), now, now),
		)

	body := `{
		"table_number": 11,
		"name": "Table 11",
		"capacity": 4,
		"location": "indoor",
		"shape": "square"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tables", bytes.NewBufferString(body))
	req = withTenantContext(req, 1, 2)
	res := httptest.NewRecorder()

	handler.CreateTable(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (%s)", res.Code, res.Body.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid json response, got error: %v", err)
	}
	if payload["success"] != true {
		t.Fatalf("expected success=true, got %v", payload["success"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableHandler_ListTablesInvalidStatus(t *testing.T) {
	handler, mock, cleanup := newTableHandlerWithMock(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tables?status=invalid", nil)
	req = withTenantContext(req, 1, 1)
	res := httptest.NewRecorder()

	handler.ListTables(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid status, got %d", res.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unexpected sql execution: %v", err)
	}
}

func TestTableHandler_UpdateStatusNotFound(t *testing.T) {
	handler, mock, cleanup := newTableHandlerWithMock(t)
	defer cleanup()

	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(99), int64(1), int64(1)).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/tables/99/status", bytes.NewBufferString(`{"status":"available"}`))
	req.SetPathValue("id", "99")
	req = withTenantContext(req, 1, 1)
	res := httptest.NewRecorder()

	handler.UpdateStatus(res, req)
	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing table, got %d", res.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableHandler_GetPublicTableInfoSuccess(t *testing.T) {
	handler, mock, cleanup := newTableHandlerWithMock(t)
	defer cleanup()

	mock.ExpectQuery("FROM restaurant_tables t").
		WithArgs(int64(5)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "table_number", "coalesce", "capacity", "location", "status", "tenant_id", "restaurant_id", "slug",
			}).AddRow(int64(5), 5, "Garden 5", 4, "indoor", "occupied", int64(1), int64(2), "demo"),
		)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/table/5/info", nil)
	req.SetPathValue("id", "5")
	res := httptest.NewRecorder()

	handler.GetPublicTableInfo(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"table_number":5`) {
		t.Fatalf("expected table_number in response body, got: %s", res.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableHandler_GetTableTenantIsolation(t *testing.T) {
	handler, mock, cleanup := newTableHandlerWithMock(t)
	defer cleanup()

	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(int64(42), int64(9), int64(9)).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tables/42", nil)
	req.SetPathValue("id", "42")
	req = withTenantContext(req, 9, 9)
	res := httptest.NewRecorder()

	handler.GetTable(res, req)
	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for cross-tenant access attempt, got %d", res.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestTableHandler_UpdatePositionsSuccess(t *testing.T) {
	handler, mock, cleanup := newTableHandlerWithMock(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs(10.0, 20.0, int64(1), int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs(30.0, 40.0, int64(2), int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body := `{"positions":[{"table_id":1,"position_x":10,"position_y":20},{"table_id":2,"position_x":30,"position_y":40}]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tables/positions", bytes.NewBufferString(body))
	req = withTenantContext(req, 1, 2)
	res := httptest.NewRecorder()

	handler.UpdatePositions(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", res.Code, res.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
