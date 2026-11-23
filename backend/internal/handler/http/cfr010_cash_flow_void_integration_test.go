//go:build integration

package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

func TestCFR010CashFlowHandlerAfterVoid(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	name := fmt.Sprintf("cfr010-http-%d", time.Now().UnixNano())
	var tenant, restaurant int
	if err := db.QueryRow(`INSERT INTO tenants (name,slug,domain) VALUES ($1,$1,$2) RETURNING id`, name, name+".example.test").Scan(&tenant); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id,name,slug,address) VALUES ($1,$2,$2,'test') RETURNING id`, tenant, name).Scan(&restaurant); err != nil {
		t.Fatal(err)
	}
	accounts := make([]int, 3)
	for i, a := range []struct{ code, label, kind, normal string }{{"1101", "Cash", "asset", "debit"}, {"2101", "Tax", "liability", "credit"}, {"4100", "Revenue", "revenue", "credit"}} {
		if err := db.QueryRow(`INSERT INTO chart_of_accounts (tenant_id,restaurant_id,account_code,account_name,account_type,normal_balance) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`, tenant, restaurant, a.code, a.label, a.kind, a.normal).Scan(&accounts[i]); err != nil {
			t.Fatal(err)
		}
	}
	entry := &domain.JournalEntry{TenantID: tenant, RestaurantID: restaurant, EntryDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), Description: "synthetic taxed cash sale", EntryType: "manual", Status: domain.JournalStatusDraft, Lines: []domain.JournalEntryLine{
		{AccountID: accounts[0], AccountCode: "1101", DebitAmount: 55, LineOrder: 1},
		{AccountID: accounts[1], AccountCode: "2101", CreditAmount: 5, LineOrder: 2},
		{AccountID: accounts[2], AccountCode: "4100", CreditAmount: 50, LineOrder: 3},
	}}
	repo := repository.NewJournalEntryRepository(db)
	id, err := repo.CreateJournalEntry(entry)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.PostEntry(tenant, restaurant, id, 1); err != nil {
		t.Fatal(err)
	}
	h := NewFinancialReportHandlerWithCache(usecase.NewFinancialReportUseCase(db, nil), service.NewInMemoryCache(), 5*time.Minute)
	assertCFR010CashFlowHTTP(t, h, tenant, restaurant, "2025-01-01", "2025-12-31", 50)
	if err := repo.VoidEntry(tenant, restaurant, id, 1, "synthetic void"); err != nil {
		t.Fatal(err)
	}
	var voidDate time.Time
	if err := db.QueryRow(`SELECT voided_at::date FROM journal_entries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, id, tenant, restaurant).Scan(&voidDate); err != nil {
		t.Fatal(err)
	}
	assertCFR010CashFlowHTTP(t, h, tenant, restaurant, "2025-01-01", "2025-12-31", 50)
	assertCFR010CashFlowHTTP(t, h, tenant, restaurant, voidDate.Format("2006-01-02"), voidDate.Format("2006-01-02"), -50)
}

func assertCFR010CashFlowHTTP(t *testing.T, h *FinancialReportHandler, tenant, restaurant int, from, to string, want float64) {
	t.Helper()
	req := withAccountingAuthContext(httptest.NewRequest(http.MethodGet, "/api/v1/reports/financial/cash-flow?date_from="+from+"&date_to="+to, nil), int64(tenant), int64(restaurant), 1)
	rec := httptest.NewRecorder()
	h.GetCashFlow(rec, req)
	var body struct {
		Data struct {
			Operating float64 `json:"operating_cash_flow"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || rec.Code != http.StatusOK || body.Data.Operating != want {
		t.Fatalf("cash flow HTTP status=%d body=%s, want operating %.2f: %v", rec.Code, rec.Body.String(), want, err)
	}
}
