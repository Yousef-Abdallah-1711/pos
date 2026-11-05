package http

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
)

func TestGetCategoryHidesDatabaseFailures(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnError(errors.New("internal database detail"))
	handler := NewCategoryHandler(repository.NewCategoryRepository(db))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/categories/7", nil)
	request.SetPathValue("id", "7")
	request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, &jwt.Claims{TenantID: 1, RestaurantID: 1}))
	response := httptest.NewRecorder()
	handler.GetCategory(response, request)
	if response.Code != http.StatusInternalServerError || strings.Contains(response.Body.String(), "internal database detail") {
		t.Fatalf("GET error response = %d %q", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations: %v", err)
	}
}

func TestGetCategoryMapsMissingRecordToNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	handler := NewCategoryHandler(repository.NewCategoryRepository(db))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/categories/7", nil)
	request.SetPathValue("id", "7")
	request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, &jwt.Claims{TenantID: 1, RestaurantID: 1}))
	response := httptest.NewRecorder()
	handler.GetCategory(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("GET missing category status = %d, want %d", response.Code, http.StatusNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations: %v", err)
	}
}
