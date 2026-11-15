//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

func TestFinancialAcceptanceLiveAPIAccountantJournals(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("POS_FINANCIAL_ACCEPTANCE_API_URL")), "/")
	if baseURL == "" {
		t.Skip("requires the dedicated disposable financial acceptance API process")
	}
	db, ctx, cancel := openFinancialAcceptanceDatabase(t)
	defer cancel()
	defer db.Close()
	fixture := seedLiveFinancialOrderFixture(t, db)
	accountant := seedLiveFinancialIdentity(t, db, fixture, "accountant", "ACCEPTANCE-ACCOUNTANT", []liveFinancialPermission{{moduleID: 15, level: "WRITE"}})
	cashier := seedLiveFinancialIdentity(t, db, fixture, "cashier", "ACCEPTANCE-CASHIER", nil)
	accountantClient := loginLiveFinancialIdentity(t, baseURL, accountant)
	cashierClient := loginLiveFinancialIdentity(t, baseURL, cashier)
	debitAccountID := liveFinancialAccountID(t, db, fixture, "5110")
	creditAccountID := liveFinancialAccountID(t, db, fixture, "1110")
	baseline := countFinancialAcceptanceJournals(t, db, ctx, fixture)
	journalID, reference := createLiveAccountantDraft(t, baseURL, db, fixture, accountantClient, debitAccountID, creditAccountID)
	postPath := financialAcceptancePath(baseURL, "/api/v1/accounting/journal-entries/%d/post", journalID)
	status, body := acceptanceJSONRequest(t, cashierClient, http.MethodPost, postPath, map[string]any{})
	if status != http.StatusForbidden {
		t.Fatalf("unauthorized cashier post returned %d, want 403: %s", status, body)
	}
	assertFinancialJournalState(t, db, ctx, fixture, journalID, "draft", 0)
	status, body = acceptanceJSONRequest(t, accountantClient, http.MethodPost, postPath, map[string]any{})
	if status != http.StatusOK {
		t.Fatalf("authorized accountant post returned %d: %s", status, body)
	}
	assertFinancialJournalState(t, db, ctx, fixture, journalID, "posted", 2)
	assertAccountantJournalGL(t, db, ctx, fixture, journalID)
	duplicateStatus, duplicateBody := acceptanceJSONRequest(t, accountantClient, http.MethodPost, postPath, map[string]any{})
	if duplicateStatus < 400 {
		t.Fatalf("duplicate accountant post returned %d, want a client error: %s", duplicateStatus, duplicateBody)
	}
	assertFinancialJournalState(t, db, ctx, fixture, journalID, "posted", 2)
	assertAccountantJournalRejected(t, baseURL, db, ctx, fixture, accountantClient, debitAccountID, creditAccountID)
	assertForeignAccountRejected(t, baseURL, db, ctx, fixture, accountantClient, debitAccountID)
	finalCount := countFinancialAcceptanceJournals(t, db, ctx, fixture)
	if finalCount != baseline+1 {
		t.Fatalf("journal count baseline/actual=%d/%d after validation failures; partial journal writes detected", baseline, finalCount)
	}
	t.Logf("financial_acceptance scenario=authenticated_accountant_journal tenant_id=%d restaurant_id=%d journal_id=%d reference=%s create_http_status=201 cashier_post_http_status=%d accountant_post_http_status=%d duplicate_post_http_status=%d journal_status=posted gl_rows=2 debit_cents=10000 credit_cents=10000 difference_cents=0 validation_rejections=3 partial_writes=0", fixture.tenantID, fixture.restaurantID, journalID, reference, statusForbidden, http.StatusOK, duplicateStatus)
}

const statusForbidden = http.StatusForbidden

func openFinancialAcceptanceDatabase(t *testing.T) (*sql.DB, context.Context, context.CancelFunc) {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable financial acceptance database: %v", err)
	}
	db.SetMaxOpenConns(8)
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	if err := db.PingContext(ctx); err != nil {
		cancel()
		_ = db.Close()
		t.Fatalf("connect disposable financial acceptance database: %v", err)
	}
	return db, ctx, cancel
}

func liveFinancialAccountID(t *testing.T, db *sql.DB, fixture liveFinancialFixture, code string) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow(`SELECT id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code=$3`, fixture.tenantID, fixture.restaurantID, code).Scan(&id); err != nil {
		t.Fatalf("read synthetic account %s: %v", code, err)
	}
	return id
}

func createLiveAccountantDraft(t *testing.T, baseURL string, db *sql.DB, fixture liveFinancialFixture, client *http.Client, debitAccountID, creditAccountID int64) (int64, string) {
	t.Helper()
	reference := fmt.Sprintf("FA-ACCOUNTANT-%d", time.Now().UnixNano())
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/accounting/journal-entries", map[string]any{
		"entry_date": time.Now().UTC().Format("2006-01-02"), "description": "Synthetic acceptance journal",
		"reference_type": "acceptance", "reference_number": reference, "entry_type": "manual",
		"lines": []map[string]any{
			{"account_id": debitAccountID, "debit_amount": 100, "credit_amount": 0, "line_order": 1, "description": "Synthetic expense"},
			{"account_id": creditAccountID, "debit_amount": 0, "credit_amount": 100, "line_order": 2, "description": "Synthetic cash"},
		},
	})
	var response struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if status != http.StatusCreated || json.Unmarshal(body, &response) != nil || response.Data.ID <= 0 {
		t.Fatalf("accountant journal create returned %d or no id: %s", status, body)
	}
	var persisted int64
	if err := db.QueryRow(`SELECT id FROM journal_entries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3 AND reference_number=$4 AND status='draft'`, response.Data.ID, fixture.tenantID, fixture.restaurantID, reference).Scan(&persisted); err != nil || persisted != response.Data.ID {
		t.Fatalf("created accountant journal is not persisted in draft state: id=%d err=%v", persisted, err)
	}
	return response.Data.ID, reference
}

func assertFinancialJournalState(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, journalID int64, expectedStatus string, expectedRows int) {
	t.Helper()
	var status string
	var rows int
	if err := db.QueryRowContext(ctx, `SELECT je.status, (SELECT COUNT(*) FROM general_ledger gl WHERE gl.tenant_id=je.tenant_id AND gl.restaurant_id=je.restaurant_id AND gl.journal_entry_id=je.id) FROM journal_entries je WHERE je.id=$1 AND je.tenant_id=$2 AND je.restaurant_id=$3`, journalID, fixture.tenantID, fixture.restaurantID).Scan(&status, &rows); err != nil || status != expectedStatus || rows != expectedRows {
		t.Fatalf("accountant journal id=%d status/GL rows=%s/%d, expected %s/%d (err=%v)", journalID, status, rows, expectedStatus, expectedRows, err)
	}
}

func assertAccountantJournalGL(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, journalID int64) {
	t.Helper()
	var debitCents, creditCents int64
	var rows int
	if err := db.QueryRowContext(ctx, `SELECT ROUND(SUM(gl.debit_amount)*100)::bigint, ROUND(SUM(gl.credit_amount)*100)::bigint, COUNT(*) FROM general_ledger gl WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.journal_entry_id=$3`, fixture.tenantID, fixture.restaurantID, journalID).Scan(&debitCents, &creditCents, &rows); err != nil {
		t.Fatalf("independently reconcile accountant journal GL: %v", err)
	}
	if debitCents != 10000 || creditCents != 10000 || debitCents-creditCents != 0 || rows != 2 {
		t.Fatalf("accountant journal GL debit/credit/rows=%d/%d/%d, expected 10000/10000/2", debitCents, creditCents, rows)
	}
}

func assertAccountantJournalRejected(t *testing.T, baseURL string, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, client *http.Client, debitAccountID, creditAccountID int64) {
	t.Helper()
	reference := fmt.Sprintf("FA-INVALID-%d", time.Now().UnixNano())
	unbalanced := map[string]any{"reference_number": reference + "-unbalanced", "entry_type": "manual", "lines": []map[string]any{
		{"account_id": debitAccountID, "debit_amount": 100, "credit_amount": 0}, {"account_id": creditAccountID, "debit_amount": 0, "credit_amount": 99},
	}}
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/accounting/journal-entries", unbalanced)
	if status != http.StatusBadRequest {
		t.Fatalf("unbalanced journal returned %d, want 400: %s", status, body)
	}
	status, body = acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/accounting/journal-entries", map[string]any{
		"reference_number": reference + "-invalid", "entry_type": "manual",
		"lines": []map[string]any{{"account_id": debitAccountID, "debit_amount": "invalid", "credit_amount": 0}},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("invalid monetary JSON returned %d, want 400: %s", status, body)
	}
	assertNoFinancialJournalReference(t, db, ctx, fixture, reference+"-unbalanced")
	assertNoFinancialJournalReference(t, db, ctx, fixture, reference+"-invalid")
}

func assertForeignAccountRejected(t *testing.T, baseURL string, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, client *http.Client, localAccountID int64) {
	t.Helper()
	foreignID := seedLiveFinancialForeignAccount(t, db)
	reference := fmt.Sprintf("FA-FOREIGN-%d", time.Now().UnixNano())
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/accounting/journal-entries", map[string]any{
		"reference_number": reference, "entry_type": "manual", "lines": []map[string]any{
			{"account_id": localAccountID, "debit_amount": 100, "credit_amount": 0},
			{"account_id": foreignID, "debit_amount": 0, "credit_amount": 100},
		},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("cross-tenant account journal returned %d, want 400: %s", status, body)
	}
	assertNoFinancialJournalReference(t, db, ctx, fixture, reference)
}

func seedLiveFinancialForeignAccount(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	key := fmt.Sprintf("financial-foreign-%d", time.Now().UnixNano())
	var tenantID, restaurantID, accountID int64
	if err := db.QueryRow(`INSERT INTO tenants (name,slug,domain) VALUES ($1,$2,$3) RETURNING id`, key, key, key+".example.test").Scan(&tenantID); err != nil {
		t.Fatalf("seed foreign acceptance tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id,name,slug,address) VALUES ($1,$2,$2,'synthetic') RETURNING id`, tenantID, key).Scan(&restaurantID); err != nil {
		t.Fatalf("seed foreign acceptance restaurant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO chart_of_accounts (tenant_id,restaurant_id,account_code,account_name,account_type,normal_balance) VALUES ($1,$2,'9991','Foreign synthetic account','asset','debit') RETURNING id`, tenantID, restaurantID).Scan(&accountID); err != nil {
		t.Fatalf("seed foreign acceptance account: %v", err)
	}
	return accountID
}

func countFinancialAcceptanceJournals(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture) int {
	t.Helper()
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='acceptance'`, fixture.tenantID, fixture.restaurantID).Scan(&count); err != nil {
		t.Fatalf("count synthetic accountant journals: %v", err)
	}
	return count
}

func assertNoFinancialJournalReference(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, reference string) {
	t.Helper()
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_number=$3`, fixture.tenantID, fixture.restaurantID, reference).Scan(&count); err != nil || count != 0 {
		t.Fatalf("rejected journal reference %s left %d persisted entries (err=%v)", reference, count, err)
	}
}
