//go:build integration

package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestGenerateBalanceSheetPostgresIntegrityStatus(t *testing.T) {
	requireDisposableFinancialTestDatabase(t)
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()

	asOf := time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)
	journalID := insertBalanceSheetFixture(t, db, asOf)
	defer deleteBalanceSheetFixture(t, db, journalID)
	report, err := NewFinancialReportUseCase(db, nil).GenerateBalanceSheet(1, 1, asOf)
	if err != nil {
		t.Fatalf("GenerateBalanceSheet returned error: %v", err)
	}
	if report.IntegrityStatus != "UNRECONCILED" || report.EquationBalanced {
		t.Fatalf("status/balanced = %s/%t, want UNRECONCILED/false", report.IntegrityStatus, report.EquationBalanced)
	}
	if report.TotalAssets != 100 || report.TotalEquity != 0 || report.EquationDifference != 100 {
		t.Fatalf("report changed source figures: assets=%v equity=%v difference=%v", report.TotalAssets, report.TotalEquity, report.EquationDifference)
	}
}

func TestGenerateBalanceSheetPostgresBackdatedJournalUsesAsOfMovementSum(t *testing.T) {
	requireDisposableFinancialTestDatabase(t)
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()

	stamp := time.Now().UnixNano()
	accountSuffix := stamp % 1_000_000_000_000
	var tenantID, restaurantID, assetID, revenueID int
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, fmt.Sprintf("Backdated report %d", stamp), fmt.Sprintf("backdated-%d", stamp), fmt.Sprintf("backdated-%d.example.test", stamp)).Scan(&tenantID); err != nil {
		t.Fatalf("create isolated tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, 'test') RETURNING id`, tenantID, "Backdated report fixture", fmt.Sprintf("backdated-%d", stamp)).Scan(&restaurantID); err != nil {
		t.Fatalf("create isolated restaurant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance) VALUES ($1,$2,$3,'Backdated cash','asset','debit') RETURNING id`, tenantID, restaurantID, fmt.Sprintf("BDA%d", accountSuffix)).Scan(&assetID); err != nil {
		t.Fatalf("create asset account: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance) VALUES ($1,$2,$3,'Backdated revenue','revenue','credit') RETURNING id`, tenantID, restaurantID, fmt.Sprintf("BDR%d", accountSuffix)).Scan(&revenueID); err != nil {
		t.Fatalf("create revenue account: %v", err)
	}
	ledger := repository.NewGeneralLedgerRepository(db, repository.NewJournalEntryRepository(db))
	post := func(referenceID int, date time.Time, amount float64) {
		t.Helper()
		entry := &domain.JournalEntry{
			TenantID: tenantID, RestaurantID: restaurantID, EntryDate: date,
			Description: "Synthetic backdated financial report fixture", ReferenceType: "test", ReferenceID: &referenceID,
			EntryType: "manual", Status: domain.JournalStatusDraft,
			Lines: []domain.JournalEntryLine{
				{AccountID: assetID, DebitAmount: amount, LineOrder: 1},
				{AccountID: revenueID, CreditAmount: amount, LineOrder: 2},
			},
		}
		if _, err := ledger.CreateAndPostAutoEntry(entry, 0); err != nil {
			t.Fatalf("post synthetic journal %d: %v", referenceID, err)
		}
	}
	lateDate := time.Date(2026, time.November, 20, 0, 0, 0, 0, time.UTC)
	earlyDate := time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)
	post(1, lateDate, 100)
	post(2, earlyDate, 50)
	post(3, earlyDate, 25)

	accountRepo := repository.NewAccountRepository(db)
	asset, err := accountRepo.GetAccountByID(tenantID, restaurantID, assetID)
	if err != nil || asset == nil {
		t.Fatalf("read posted asset account: account=%+v err=%v", asset, err)
	}
	asset.NormalBalance = "credit"
	if err := accountRepo.UpdateAccount(asset); err == nil {
		t.Fatal("changing normal_balance after ledger posting unexpectedly succeeded")
	}
	asset, err = accountRepo.GetAccountByID(tenantID, restaurantID, assetID)
	if err != nil || asset == nil || asset.NormalBalance != "debit" {
		t.Fatalf("rejected normal_balance change altered source account: account=%+v err=%v", asset, err)
	}

	report, err := NewFinancialReportUseCase(db, nil).GenerateBalanceSheet(tenantID, restaurantID, earlyDate)
	if err != nil {
		t.Fatalf("GenerateBalanceSheet returned error: %v", err)
	}
	assertUsecaseBalanceCents(t, "Balance Sheet assets as of backdated date", 7500, report.TotalAssets)

	latestDate := time.Date(2026, time.November, 30, 0, 0, 0, 0, time.UTC)
	post(4, latestDate, 10)
	latestReport, err := NewFinancialReportUseCase(db, nil).GenerateBalanceSheet(tenantID, restaurantID, latestDate)
	if err != nil {
		t.Fatalf("GenerateBalanceSheet returned error for latest date: %v", err)
	}
	assertUsecaseBalanceCents(t, "Balance Sheet assets after backdated and subsequent postings", 18500, latestReport.TotalAssets)
	currentBalance, err := ledger.GetAccountBalance(tenantID, restaurantID, assetID, nil)
	if err != nil {
		t.Fatalf("read movement balance after backdated posting: %v", err)
	}
	assertUsecaseBalanceCents(t, "current movement balance", 18500, currentBalance)
	asset, err = accountRepo.GetAccountByID(tenantID, restaurantID, assetID)
	if err != nil || asset == nil {
		t.Fatalf("read current account balance: account=%+v err=%v", asset, err)
	}
	assertUsecaseBalanceCents(t, "stored current account balance", 18500, asset.CurrentBalance)
	ledgerRows, err := ledger.GetLedgerByAccount(tenantID, restaurantID, assetID, nil, nil)
	if err != nil || len(ledgerRows) != 4 {
		t.Fatalf("read chronological ledger rows: rows=%+v err=%v, want 4 rows", ledgerRows, err)
	}
	for index, expected := range []int64{18500, 17500, 7500, 5000} {
		assertUsecaseBalanceCents(t, fmt.Sprintf("ledger row %d running balance", index+1), expected, ledgerRows[index].RunningBalance)
	}
}

func TestFinancialReportsPostgresReconcileToIndependentLedgerOracle(t *testing.T) {
	requireDisposableFinancialTestDatabase(t)
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()

	stamp := time.Now().UnixNano()
	var tenantID, restaurantID, assetID, liabilityID, revenueID int
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, fmt.Sprintf("Reconciliation oracle %d", stamp), fmt.Sprintf("reconcile-%d", stamp), fmt.Sprintf("reconcile-%d.example.test", stamp)).Scan(&tenantID); err != nil {
		t.Fatalf("create isolated tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, 'test') RETURNING id`, tenantID, "Reconciliation oracle fixture", fmt.Sprintf("reconcile-%d", stamp)).Scan(&restaurantID); err != nil {
		t.Fatalf("create isolated restaurant: %v", err)
	}
	accountSuffix := stamp % 1_000_000_000_000
	for _, account := range []struct {
		code, name, kind, normal string
		destination              *int
	}{
		{fmt.Sprintf("ROA%d", accountSuffix), "Oracle cash", "asset", "debit", &assetID},
		{fmt.Sprintf("ROL%d", accountSuffix), "Oracle payable", "liability", "credit", &liabilityID},
		{fmt.Sprintf("ROR%d", accountSuffix), "Oracle revenue", "revenue", "credit", &revenueID},
	} {
		if err := db.QueryRow(`INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`, tenantID, restaurantID, account.code, account.name, account.kind, account.normal).Scan(account.destination); err != nil {
			t.Fatalf("create %s account: %v", account.kind, err)
		}
	}

	entryDate := time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC)
	referenceID := int(stamp % 1_000_000_000)
	entry := &domain.JournalEntry{
		TenantID: tenantID, RestaurantID: restaurantID, EntryDate: entryDate,
		Description: "Synthetic reconciliation oracle fixture", ReferenceType: "test", ReferenceID: &referenceID,
		EntryType: "manual", Status: domain.JournalStatusDraft,
		Lines: []domain.JournalEntryLine{
			{AccountID: assetID, DebitAmount: 55, LineOrder: 1},
			{AccountID: liabilityID, CreditAmount: 5, LineOrder: 2},
			{AccountID: revenueID, CreditAmount: 50, LineOrder: 3},
		},
	}
	ledger := repository.NewGeneralLedgerRepository(db, repository.NewJournalEntryRepository(db))
	journalID, err := ledger.CreateAndPostAutoEntry(entry, 0)
	if err != nil {
		t.Fatalf("post synthetic reconciliation journal: %v", err)
	}

	// Oracle calculation reads source posted journal lines directly; it does not
	// reuse the report query, GL running-balance cache, or report mapping code.
	var sourceDebits, sourceCredits float64
	if err := db.QueryRow(`
		SELECT COALESCE(SUM(jel.debit_amount), 0), COALESCE(SUM(jel.credit_amount), 0)
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		WHERE je.id = $1 AND je.tenant_id = $2 AND je.restaurant_id = $3 AND je.status = 'posted'
	`, journalID, tenantID, restaurantID).Scan(&sourceDebits, &sourceCredits); err != nil {
		t.Fatalf("sum source journal lines: %v", err)
	}
	if sourceDebits != 55 || sourceCredits != 55 {
		t.Fatalf("source journal oracle debits/credits = %.2f/%.2f, want 55/55 currency units", sourceDebits, sourceCredits)
	}

	type accountTotals struct{ debit, credit float64 }
	oracle := map[string]accountTotals{}
	rows, err := db.Query(`
		SELECT coa.account_type, COALESCE(SUM(gl.debit_amount), 0), COALESCE(SUM(gl.credit_amount), 0)
		FROM chart_of_accounts coa
		LEFT JOIN general_ledger gl ON gl.account_id = coa.id AND gl.tenant_id = $1 AND gl.restaurant_id = $2 AND gl.entry_date <= $3
		WHERE coa.tenant_id = $1 AND coa.restaurant_id = $2 AND coa.is_active = TRUE AND coa.is_header = FALSE
		GROUP BY coa.account_type
	`, tenantID, restaurantID, entryDate)
	if err != nil {
		t.Fatalf("sum independent ledger movements: %v", err)
	}
	for rows.Next() {
		var kind string
		var totals accountTotals
		if err := rows.Scan(&kind, &totals.debit, &totals.credit); err != nil {
			rows.Close()
			t.Fatalf("scan independent account totals: %v", err)
		}
		oracle[kind] = totals
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		t.Fatalf("read independent account totals: %v", err)
	}
	rows.Close()
	assetCents := int64(math.Round((oracle["asset"].debit - oracle["asset"].credit) * 100))
	liabilityCents := int64(math.Round((oracle["liability"].credit - oracle["liability"].debit) * 100))
	equityCents := int64(math.Round((oracle["equity"].credit - oracle["equity"].debit) * 100))
	revenueCents := int64(math.Round((oracle["revenue"].credit - oracle["revenue"].debit) * 100))
	expenseCents := int64(math.Round((oracle["expense"].debit - oracle["expense"].credit) * 100))
	if assetCents != 5500 || liabilityCents != 500 || equityCents != 0 || revenueCents != 5000 || expenseCents != 0 {
		t.Fatalf("independent ledger oracle cents assets=%d liabilities=%d equity=%d revenue=%d expense=%d", assetCents, liabilityCents, equityCents, revenueCents, expenseCents)
	}

	usecase := NewFinancialReportUseCase(db, ledger)
	trialBalance, err := usecase.GenerateTrialBalance(tenantID, restaurantID, entryDate)
	if err != nil {
		t.Fatalf("generate trial balance: %v", err)
	}
	if !trialBalance.IsBalanced || math.Round(trialBalance.TotalDebits*100) != 5500 || math.Round(trialBalance.TotalCredits*100) != 5500 {
		t.Fatalf("trial balance diverges from source oracle: %+v", trialBalance)
	}
	profitLoss, err := usecase.GenerateProfitLossWithComparison(tenantID, restaurantID, entryDate, entryDate, entryDate.AddDate(0, 0, -1), entryDate.AddDate(0, 0, -1))
	if err != nil {
		t.Fatalf("generate profit and loss: %v", err)
	}
	assertUsecaseBalanceCents(t, "P&L revenue from independent oracle", revenueCents, profitLoss.TotalRevenue)
	assertUsecaseBalanceCents(t, "P&L expense from independent oracle", expenseCents, profitLoss.TotalExpenses)
	assertUsecaseBalanceCents(t, "P&L net income from independent oracle", revenueCents-expenseCents, profitLoss.NetProfit)
	balanceSheet, err := usecase.GenerateBalanceSheet(tenantID, restaurantID, entryDate)
	if err != nil {
		t.Fatalf("generate balance sheet: %v", err)
	}
	assertUsecaseBalanceCents(t, "Balance Sheet assets from independent oracle", assetCents, balanceSheet.TotalAssets)
	assertUsecaseBalanceCents(t, "Balance Sheet liabilities from independent oracle", liabilityCents, balanceSheet.TotalLiabilities)
	assertUsecaseBalanceCents(t, "Balance Sheet equity from independent oracle", equityCents, balanceSheet.TotalEquity)
	assertUsecaseBalanceCents(t, "unclosed-earnings residual from independent oracle", revenueCents-expenseCents, balanceSheet.EquationDifference)
	if balanceSheet.IntegrityStatus != "UNRECONCILED" || balanceSheet.EquationBalanced {
		t.Fatalf("Balance Sheet hid known unclosed-earnings difference: %+v", balanceSheet)
	}
}

func assertUsecaseBalanceCents(t *testing.T, label string, expected int64, observedAmount float64) {
	t.Helper()
	observed := int64(math.Round(observedAmount * 100))
	difference := observed - expected
	t.Logf("financial_reconciliation scenario=%s expected_cents=%d observed_cents=%d difference_cents=%d", label, expected, observed, difference)
	if difference != 0 {
		t.Fatalf("%s reconciliation: expected=%d cents observed=%d cents difference=%d cents", label, expected, observed, difference)
	}
}

func requireDisposableFinancialTestDatabase(t *testing.T) {
	t.Helper()
	if os.Getenv("SERVIO_DISPOSABLE_FINANCIAL_TEST_DB") != "1" {
		t.Skip("this test posts immutable synthetic ledger history; set SERVIO_DISPOSABLE_FINANCIAL_TEST_DB=1 only for a disposable PostgreSQL database")
	}
}

func insertBalanceSheetFixture(t *testing.T, db *sql.DB, entryDate time.Time) int64 {
	t.Helper()
	ctx := context.Background()
	entryNumber := fmt.Sprintf("P00-%d", time.Now().UnixNano()%1_000_000_000_000)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin fixture transaction: %v", err)
	}
	defer func() { _ = tx.Rollback() }()
	var journalID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, description, entry_type, status)
		VALUES (1, 1, $1, $2, 'Phase 00 integrity fixture', 'manual', 'draft') RETURNING id
	`, entryNumber, entryDate).Scan(&journalID)
	if err != nil {
		t.Fatalf("insert fixture journal: %v", err)
	}
	fixtureLines := []struct {
		code, description      string
		debit, credit, balance float64
	}{
		{code: "1110", description: "Cash on Hand", debit: 100, balance: 100},
		{code: "4100", description: "Food Sales", credit: 100, balance: -100},
	}
	accountIDs := make([]int64, len(fixtureLines))
	lineIDs := make([]int64, len(fixtureLines))
	for index, fixtureLine := range fixtureLines {
		if err := tx.QueryRowContext(ctx, `SELECT id FROM chart_of_accounts WHERE tenant_id = $1 AND restaurant_id = $2 AND account_code = $3`, 1, 1, fixtureLine.code).Scan(&accountIDs[index]); err != nil {
			t.Fatalf("find seeded account %s: %v", fixtureLine.code, err)
		}
		err = tx.QueryRowContext(ctx, `
			INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, credit_amount, description, line_order)
			VALUES ($1, 1, $2, $3, $4, $5, $6, $7) RETURNING id
		`, journalID, accountIDs[index], fixtureLine.code, fixtureLine.debit, fixtureLine.credit, fixtureLine.description, index+1).Scan(&lineIDs[index])
		if err != nil {
			t.Fatalf("insert fixture journal line %s: %v", fixtureLine.code, err)
		}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE journal_entries SET status = 'posted', posted_at = NOW(), posted_by = 1 WHERE id = $1`, journalID); err != nil {
		t.Fatalf("transition fixture journal to posted: %v", err)
	}
	for index, fixtureLine := range fixtureLines {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO general_ledger (tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance)
			VALUES (1, 1, $1, $2, $3, $4, '2026-11', $5, $6, $7)
		`, accountIDs[index], journalID, lineIDs[index], entryDate, fixtureLine.debit, fixtureLine.credit, fixtureLine.balance)
		if err != nil {
			t.Fatalf("insert fixture ledger row %s: %v", fixtureLine.code, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit fixture: %v", err)
	}
	return journalID
}

func deleteBalanceSheetFixture(t *testing.T, db *sql.DB, journalID int64) {
	t.Helper()
	// The fixture is posted accounting history, which CFR-002 makes immutable.
	// The dedicated disposable PostgreSQL database is discarded after the run.
}
