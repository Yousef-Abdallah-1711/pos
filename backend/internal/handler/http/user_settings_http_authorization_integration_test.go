package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestUserSettingsHTTPAllowsSameUserSettingsAndProfileOperations(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	userID := insertSettingsHTTPUser(t, db, fixture.tenantA, fixture.restaurantA1, "settings-owner")
	if _, err := db.Exec(`INSERT INTO user_settings (user_id, language) VALUES ($1, 'en')`, userID); err != nil {
		t.Fatalf("insert settings fixture: %v", err)
	}
	initialPasswordHash, err := bcrypt.GenerateFromPassword([]byte("current-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash current password fixture: %v", err)
	}
	if _, err := db.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, string(initialPasswordHash), userID); err != nil {
		t.Fatalf("set current password fixture: %v", err)
	}
	handler, path, token := sameUserSettingsHTTPHandler(t, db, fixture, userID)
	assertSameUserCanReadSettings(t, handler, path, token)
	assertSameUserCanUpdateLanguage(t, db, handler, path, token, userID)
	assertSameUserCanUpdateThemeAndColors(t, db, handler, path, token, userID)
	assertSameUserCanChangePassword(t, db, handler, path, token, userID)
	assertSameUserCanReadAndUpdateProfile(t, db, handler, path, token, userID)
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
	mux.HandleFunc("PUT "+userPath+"/settings/theme", handler.UpdateTheme)
	mux.HandleFunc("PUT "+userPath+"/settings/theme-colors", handler.UpdateColors)
	mux.HandleFunc("POST "+userPath+"/settings/change-password", handler.ChangePassword)
	mux.HandleFunc("GET "+userPath+"/profile", handler.GetUserProfile)
	mux.HandleFunc("PUT "+userPath+"/profile", handler.UpdateProfile)
	secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
		middleware.SQLRestaurantVerifier{DB: db},
	)(mux))
	return secured, fmt.Sprintf("/api/v1/users/%d/settings", userID), token
}

func assertSameUserCanUpdateThemeAndColors(t *testing.T, db *sql.DB, secured http.Handler, path, token string, userID int64) {
	t.Helper()
	assertAuthenticatedSettingsStatus(t, secured, http.MethodPut, path+"/theme", token, `{"theme":"dark"}`)
	assertAuthenticatedSettingsStatus(t, secured, http.MethodPut, path+"/theme-colors", token, `{"primary_color":"#112233","secondary_color":"#445566","accent_color":"#778899"}`)
	var theme, primary, secondary, accent string
	err := db.QueryRow(`SELECT theme, primary_color, secondary_color, accent_color FROM user_settings WHERE user_id = $1`, userID).Scan(&theme, &primary, &secondary, &accent)
	if err != nil {
		t.Fatalf("read updated theme and colors: %v", err)
	}
	if theme != "dark" || primary != "#112233" || secondary != "#445566" || accent != "#778899" {
		t.Fatalf("same-user settings stored theme/colors %q %q %q %q", theme, primary, secondary, accent)
	}
}

func assertSameUserCanChangePassword(t *testing.T, db *sql.DB, secured http.Handler, path, token string, userID int64) {
	t.Helper()
	assertAuthenticatedSettingsStatus(t, secured, http.MethodPost, path+"/change-password", token, `{"current_password":"current-password","new_password":"updated-password","confirm_password":"updated-password"}`)
	var passwordHash string
	if err := db.QueryRow(`SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&passwordHash); err != nil {
		t.Fatalf("read changed password hash: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("updated-password")); err != nil {
		t.Fatalf("same-user password change did not persist the new password: %v", err)
	}
	var historyCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM password_change_history WHERE user_id = $1`, userID).Scan(&historyCount); err != nil {
		t.Fatalf("count password change history: %v", err)
	}
	if historyCount != 1 {
		t.Fatalf("same-user password change wrote %d history rows, want 1", historyCount)
	}
}

func assertSameUserCanReadAndUpdateProfile(t *testing.T, db *sql.DB, secured http.Handler, path, token string, userID int64) {
	t.Helper()
	const avatarURL = "https://example.test/settings-owner.png"
	profilePath := strings.TrimSuffix(path, "/settings") + "/profile"
	update := authenticatedSettingsRequest(http.MethodPut, profilePath, token, `{"name":"Updated Settings Owner","avatar_url":"`+avatarURL+`"}`)
	updateRecorder := httptest.NewRecorder()
	secured.ServeHTTP(updateRecorder, update)
	if updateRecorder.Code != http.StatusOK || !strings.Contains(updateRecorder.Body.String(), "Updated Settings Owner") || !strings.Contains(updateRecorder.Body.String(), avatarURL) {
		t.Fatalf("same-user profile update returned %d: %s", updateRecorder.Code, updateRecorder.Body.String())
	}
	read := authenticatedSettingsRequest(http.MethodGet, profilePath, token, "")
	readRecorder := httptest.NewRecorder()
	secured.ServeHTTP(readRecorder, read)
	if readRecorder.Code != http.StatusOK || !strings.Contains(readRecorder.Body.String(), "Updated Settings Owner") {
		t.Fatalf("same-user profile read returned %d: %s", readRecorder.Code, readRecorder.Body.String())
	}
	var name, storedAvatar string
	if err := db.QueryRow(`SELECT name, avatar_url FROM users WHERE id = $1`, userID).Scan(&name, &storedAvatar); err != nil {
		t.Fatalf("read updated profile state: %v", err)
	}
	if name != "Updated Settings Owner" || storedAvatar != avatarURL {
		t.Fatalf("same-user profile stored name/avatar %q %q", name, storedAvatar)
	}
}

func assertAuthenticatedSettingsStatus(t *testing.T, secured http.Handler, method, path, token, body string) {
	t.Helper()
	request := authenticatedSettingsRequest(method, path, token, body)
	recorder := httptest.NewRecorder()
	secured.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("same-user %s %s returned %d: %s", method, path, recorder.Code, recorder.Body.String())
	}
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
