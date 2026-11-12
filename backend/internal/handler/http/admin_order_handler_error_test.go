package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func TestAdminOrderListHidesDatabaseDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT COUNT\\(DISTINCT o.id\\)").WillReturnError(errors.New("internal database detail"))
	handler := NewAdminOrderHandler(usecase.NewOrderUseCase(repository.NewOrderRepository(db), nil), nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orders", nil)
	ctx := context.WithValue(request.Context(), middleware.TenantContextKey, int64(1))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(1))
	response := httptest.NewRecorder()
	handler.ListOrders(response, request.WithContext(ctx))
	if response.Code != http.StatusInternalServerError || strings.Contains(response.Body.String(), "internal database detail") {
		t.Fatalf("list error response = %d %q", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations: %v", err)
	}
}
