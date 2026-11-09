package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestBankAccountHTTPRejectsForeignTenantAndRestaurantMutations(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db := openProductHTTPTestDB(t)
	fixture := seedProductHTTPFixture(t, db)
	targets := map[string]struct {
		tenantID     int
		restaurantID int
	}{
		"foreign tenant":     {fixture.tenantB, fixture.restaurantB1},
		"sibling restaurant": {fixture.tenantA, fixture.restaurantA2},
	}
	accountIDs := make(map[string]int64, len(targets))
	for name, target := range targets {
		accountIDs[name] = insertBankAccountIsolationFixture(t, db, target.tenantID, target.restaurantID, name)
	}

	tokens, err := jwt.NewTokenService("ws02-bank-account-test-secret-at-least-32-bytes", "5m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	restaurantID := fixture.restaurantA1
	token, err := tokens.GenerateToken(1, fixture.tenantA, &restaurantID, "accounting-owner@example.test", "owner")
	if err != nil {
		t.Fatalf("generate authenticated token: %v", err)
	}
	handler := NewBankAccountHandler(repository.NewBankAccountRepository(db))

	for name, target := range targets {
		t.Run(name, func(t *testing.T) {
			accountID := accountIDs[name]
			for _, attempt := range []struct {
				name       string
				method     string
				body       string
				handler    http.HandlerFunc
				wantStatus int
			}{
				{name: "read", method: http.MethodGet, handler: handler.GetBankAccount, wantStatus: http.StatusNotFound},
				{name: "update", method: http.MethodPut, body: `{"account_name":"tampered","bank_name":"tampered","account_number":"tampered","opening_balance":999}`, handler: handler.UpdateBankAccount, wantStatus: http.StatusNotFound},
				{name: "delete", method: http.MethodDelete, handler: handler.DeleteBankAccount, wantStatus: http.StatusNotFound},
			} {
				t.Run(attempt.name, func(t *testing.T) {
					request := httptest.NewRequest(attempt.method, fmt.Sprintf("/api/v1/accounting/bank-accounts/%d", accountID), strings.NewReader(attempt.body))
					request.SetPathValue("id", strconv.FormatInt(accountID, 10))
					request.Header.Set("Authorization", "Bearer "+token)
					request.Header.Set("X-Tenant-ID", strconv.Itoa(target.tenantID))
					request.Header.Set("X-Restaurant-ID", strconv.Itoa(target.restaurantID))
					secured := middleware.AuthMiddleware(tokens)(middleware.NewTenantContextMiddleware(
						middleware.SQLRestaurantVerifier{DB: db},
					)(attempt.handler))
					recorder := httptest.NewRecorder()
					secured.ServeHTTP(recorder, request)
					if recorder.Code != attempt.wantStatus {
						t.Fatalf("foreign bank account request returned %d, want %d: %s", recorder.Code, attempt.wantStatus, recorder.Body.String())
					}
				})
			}

			var accountName, accountNumber string
			if err := db.QueryRow(`SELECT account_name, account_number FROM bank_accounts_ledger WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, accountID, target.tenantID, target.restaurantID).Scan(&accountName, &accountNumber); err != nil {
				t.Fatalf("read target bank account after denied requests: %v", err)
			}
			if accountName != "Scoped account "+name || accountNumber != "WS02-"+name {
				t.Fatalf("denied request changed account to name=%q number=%q", accountName, accountNumber)
			}
		})
	}
}

func insertBankAccountIsolationFixture(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(`
		INSERT INTO bank_accounts_ledger (
			tenant_id, restaurant_id, account_name, bank_name, account_number, currency,
			account_type, opening_balance, current_balance, opening_date
		) VALUES ($1, $2, $3, 'Fixture Bank', $4, 'EGP', 'checking', 25, 25, $5)
		RETURNING id
	`, tenantID, restaurantID, "Scoped account "+suffix, "WS02-"+suffix, time.Now().UTC().Format("2006-01-02")).Scan(&id)
	if err != nil {
		t.Fatalf("insert bank account isolation fixture: %v", err)
	}
	return id
}
