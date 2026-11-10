package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
)

func TestRegisterReturnsFieldValidationErrorsBeforeDatabaseAccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	handler := NewAuthHandler(usecase.NewAuthUseCase(repository.NewAuthRepository(db), nil), nil)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"name":"J","email":"bad","password":"123456","restaurant_name":"Cafe","restaurant_slug":"cafe","phone":"6505550123"}`))
	response := httptest.NewRecorder()
	handler.Register(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var body authErrorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "VALIDATION_ERROR" || body.Fields["name"] == "" || body.Fields["email"] == "" || body.Fields["phone"] == "" {
		t.Fatalf("unexpected validation response: %#v", body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("validation should happen before SQL: %v", err)
	}
}

func TestRegisterHidesUnexpectedDatabaseErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM tenants WHERE slug = \$1\)`).
		WithArgs("safe-slug").WillReturnError(errors.New("pq: private table and column details"))
	handler := NewAuthHandler(usecase.NewAuthUseCase(repository.NewAuthRepository(db), nil), nil)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"name":"Jane Doe","email":"jane@example.com","password":"123456","restaurant_name":"Cafe North","restaurant_slug":"safe-slug"}`))
	response := httptest.NewRecorder()
	handler.Register(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "private table") || strings.Contains(response.Body.String(), "column details") {
		t.Fatalf("database details leaked: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"code":"INTERNAL_ERROR"`) {
		t.Fatalf("missing stable error contract: %s", response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAuthDecoderRejectsOversizedBody(t *testing.T) {
	body := `{"name":"` + strings.Repeat("x", maxAuthRequestBytes) + `"}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	response := httptest.NewRecorder()
	var payload map[string]interface{}
	if err := decodeAuthJSON(response, request, &payload); err == nil {
		t.Fatal("expected oversized request error")
	} else {
		respondAuthDecodeError(response, request, err)
	}
	if response.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
