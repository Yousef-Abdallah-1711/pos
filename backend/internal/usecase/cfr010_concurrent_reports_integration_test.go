//go:build integration

package usecase

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestCFR010ConcurrentPostAndBalanceSheetRead(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(8)
	f := createCFR010Fixture(t, db)
	date := cfr010Date("2026-04-30")
	start := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		<-start
		repo := repository.NewJournalEntryRepository(db)
		for i := 0; i < 20; i++ {
			entry := &domain.JournalEntry{TenantID: f.tenantID, RestaurantID: f.restaurantID, EntryDate: date, Description: fmt.Sprintf("concurrent synthetic %d", i), EntryType: "manual", Status: domain.JournalStatusDraft, Lines: []domain.JournalEntryLine{
				{AccountID: f.accounts["cash"], AccountCode: f.accountCodes["cash"], DebitAmount: 1, LineOrder: 1},
				{AccountID: f.accounts["capital"], AccountCode: f.accountCodes["capital"], CreditAmount: 1, LineOrder: 2},
			}}
			id, err := repo.CreateJournalEntry(entry)
			if err == nil {
				err = repo.PostEntry(f.tenantID, f.restaurantID, id, 1)
			}
			if err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	close(start)
	uc := NewFinancialReportUseCase(db, nil)
	reads := 0
	for {
		report, err := uc.GenerateBalanceSheet(f.tenantID, f.restaurantID, date)
		if err != nil {
			t.Fatal(err)
		}
		if report.TotalAssets != report.TotalEquity || !report.EquationBalanced || report.TotalAssets < 0 || report.TotalAssets > 20 {
			t.Fatalf("partial or invalid snapshot: %+v", report)
		}
		reads++
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
			final, err := uc.GenerateBalanceSheet(f.tenantID, f.restaurantID, date)
			if err != nil || final.TotalAssets != 20 || final.TotalEquity != 20 {
				t.Fatalf("final report = %+v, want 20/20: %v", final, err)
			}
			var debitCents, creditCents int64
			err = db.QueryRow(`SELECT (COALESCE(SUM(debit_amount),0)*100)::bigint, (COALESCE(SUM(credit_amount),0)*100)::bigint FROM general_ledger WHERE tenant_id=$1 AND restaurant_id=$2 AND entry_date <= $3`, f.tenantID, f.restaurantID, date).Scan(&debitCents, &creditCents)
			if err != nil || debitCents != 2000 || creditCents != 2000 {
				t.Fatalf("GL debit/credit = %d/%d cents, want 2000/2000: %v", debitCents, creditCents, err)
			}
			t.Logf("%d concurrent reads; final assets/equity=2000/2000 cents; GL debit/credit=2000/2000 cents", reads)
			return
		default:
			if reads > 10000 {
				t.Fatal("writer did not finish within read bound")
			}
			time.Sleep(time.Millisecond)
		}
	}
}
