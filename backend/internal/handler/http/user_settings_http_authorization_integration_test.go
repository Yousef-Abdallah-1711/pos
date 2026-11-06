package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestUserSettingsHTTPAllowsSameUserReadAndLanguageUpdate(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	userID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "settings-owner")
	if _, err := db.Exec(`INSERT INTO user_settings (user_id, language) VALUES ($1, 'en')`, userID); err != nil {
		t.Fatalf("insert settings fixture: %v", err)
	}
	handler, path, token := sameUserSettingsHTTPHandler(t, db, fixture, userID)
	assertSameUserCanReadSettings(t, handler, path, token)
	assertSameUserCanUpdateLanguage(t, db, handler, path, token, userID)
}

func sameUserSettingsHTTPHandler(t *testing.T, db *sql.DB, fixture productHTTPFixture, userID int64) (http.Handler, string, string) {
	t.Helper()
	tokens, err := jwt.NewTokenService("ws02-settings-owner-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(int(userID), fixture.tenantA, &restaurantID, "settings-owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate owner token: %v", err)
	}
	handler := NewUserSettingsHandler(repository.NewUserSettingsRepository(db), repository.NewUserRepository(db))
	mux := http.NewServeMux()
	userPath := "/api/v1/users/{userId}"
	mux.HandleFunc("GET "+userPath+"/settings", handler.GetUserSettings)
	mux.HandleFunc("PUT "+userPath+"/settings/language", handler.UpdateLanguage)
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(mux))
	return secured, fmt.Sprintf("/api/v1/users/%d/settings", userID), token
}

func assertSameUserCanReadSettings(t *testing.T, secured http.Handler, path, token string) {
	t.Helper()
	read := authenticatedSettingsRequest(http.MethodGet, path, token, "")
	readRecorder := httptest.NewRecorder()
	secured.ServeHTTP(readRecorder, read)
	if readRecorder.Code != http.StatusOK || !strings.Contains(readRecorder.Body.String(), `"language":"en"`) {
		t.Fatalf("same-user settings read returned %d: %s", readRecorder.Code, readRecorder.Body.String())
	}
}

func assertSameUserCanUpdateLanguage(t *testing.T, db *sql.DB, secured http.Handler, path, token string, userID int64) {
	t.Helper()
	update := authenticatedSettingsRequest(http.MethodPut, path+"/language", token, `{"language":"ar"}`)
	updateRecorder := httptest.NewRecorder()
	secured.ServeHTTP(updateRecorder, update)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("same-user language update returned %d: %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	var language string
	if err := db.QueryRow(`SELECT language FROM user_settings WHERE user_id = $1`, userID).Scan(&language); err != nil {
		t.Fatalf("read updated language: %v", err)
	}
	if language != "ar" {
		t.Fatalf("same-user update stored language %q, want ar", language)
	}
	var auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM user_settings_audit_log WHERE user_id = $1`, userID).Scan(&auditCount); err != nil {
		t.Fatalf("count settings audit rows: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("same-user update created %d audit rows, want 1", auditCount)
	}
}

func authenticatedSettingsRequest(method, path, token, body string) *http.Request {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	return request
}
